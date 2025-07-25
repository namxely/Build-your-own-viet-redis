package core

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/namxely/Build-your-own-viet-redis/cluster/raft"
	"github.com/namxely/Build-your-own-viet-redis/interface/redis"
	"github.com/namxely/Build-your-own-viet-redis/lib/logger"
	"github.com/namxely/Build-your-own-viet-redis/lib/utils"
	"github.com/namxely/Build-your-own-viet-redis/redis/protocol"
)

/* 

**Rebalance Procedure**
1. Invoke `triggerMigrationTask` on cluster Leader to start a migration task
2. Leader propose EventStartMigrate to raft group then send startMigrationCommand to TargetNode. (at `triggerMigrationTask`)

3. Target Node runs `doImports` after receiving startMigrationCommand
4. Target Node send exportCommand to Source Node.

6. Source Node get migrating task from raft (at `execExport`)
7. SourceNode set task into slotManager to start recording dirty keys during migration. (at `injectInsertCallback`)
8. Source Node dump old data to Target Node

9. Target node send migrationDoneCommand to Source Node. (at `doImports`)
10. Source Node runs `execFinishExport`, lock slots to stop writing
11. Source Node send dirty keys to Target Node

12. Source Node send migrationChangeRouteCommand to Leader
13. Leader porposes EventFinishMigrate to raft and waits Source Node and Target Node receives this entry(at `execMigrationChangeRoute`)
14. Source Node finish exporting, unlock slots, clean data
15. Target Node finish importing, unlock slots, start serve
*/

const joinClusterCommand = "cluster.join"
const migrationChangeRouteCommand = "cluster.migration.changeroute"

func init() {
	RegisterCmd(joinClusterCommand, execJoin)
	RegisterCmd(migrationChangeRouteCommand, execMigrationChangeRoute)
}

// execJoin handles cluster-join command as raft leader
// format: cluster-join redisAddress(advertised), raftAddress, masterId
func execJoin(cluster *Cluster, c redis.Connection, cmdLine CmdLine) redis.Reply {
	if len(cmdLine) < 3 {
		return protocol.MakeArgNumErrReply(joinClusterCommand)
	}
	state := cluster.raftNode.State()
	if state != raft.Leader {
		// I am not leader, forward request to leader
		leaderConn, err := cluster.BorrowLeaderClient()
		if err != nil {
			return protocol.MakeErrReply(err.Error())
		}
		defer cluster.connections.ReturnPeerClient(leaderConn)
		return leaderConn.Send(cmdLine)
	}
	// self node is leader
	redisAddr := string(cmdLine[1])
	raftAddr := string(cmdLine[2])
	err := cluster.raftNode.AddToRaft(redisAddr, raftAddr)
	if err != nil {
		return protocol.MakeErrReply(err.Error())
	}
	master := ""
	if len(cmdLine) == 4 {
		master = string(cmdLine[3])
	}
	_, err = cluster.raftNode.Propose(&raft.LogEntry{
		Event: raft.EventJoin,
		JoinTask: &raft.JoinTask{
			NodeId: redisAddr,
			Master: master,
		},
	})
	if err != nil {
		// todo: remove the node from raft
		return protocol.MakeErrReply(err.Error())
	}

	// join sucees, rebalance node
	return protocol.MakeOkReply()
}

type rebalanceManager struct {
	mu *sync.Mutex
}

func newRebalanceManager() *rebalanceManager {
	return &rebalanceManager{
		mu: &sync.Mutex{},
	}
}

func (cluster *Cluster) doRebalance() {
	cluster.rebalanceManger.mu.Lock()
	defer cluster.rebalanceManger.mu.Unlock()
	pendingTasks, err := cluster.makeRebalancePlan()
	if err != nil {
		logger.Errorf("makeRebalancePlan err: %v", err)
		return
	}
	if len(pendingTasks) == 0 {
		return
	}
	logger.Infof("rebalance plan generated, contains %d tasks", len(pendingTasks))

	for _, task := range pendingTasks {
		err := cluster.triggerMigrationTask(task)
		if err != nil {
			logger.Errorf("triggerMigrationTask err: %v", err)
		} else {
			logger.Infof("triggerMigrationTask %s success", task.ID)
		}
	}

}

// triggerRebalanceTask start a rebalance task
// only leader can do
func (cluster *Cluster) triggerMigrationTask(task *raft.MigratingTask) error {
	// propose migration
	_, err := cluster.raftNode.Propose(&raft.LogEntry{
		Event:         raft.EventStartMigrate,
		MigratingTask: task,
	})
	if err != nil {
		return fmt.Errorf("propose EventStartMigrate  %s failed: %v", task.ID, err)
	}
	logger.Infof("propose EventStartMigrate %s success", task.ID)

	cmdLine := utils.ToCmdLine(startMigrationCommand, task.ID)
	targetNodeConn, err := cluster.connections.BorrowPeerClient(task.TargetNode)
	if err != nil {
		return err
	}
	defer cluster.connections.ReturnPeerClient(targetNodeConn)
	reply := targetNodeConn.Send(cmdLine)
	if protocol.IsOKReply(reply) {
		return nil
	}
	return protocol.MakeErrReply("")
}

func (cluster *Cluster) makeRebalancePlan() ([]*raft.MigratingTask, error) {

	var migratings []*raft.MigratingTask
	cluster.raftNode.FSM.WithReadLock(func(fsm *raft.FSM) {
		avgSlot := int(math.Ceil(float64(SlotCount) / float64(len(fsm.MasterSlaves))))
		var exportingNodes []string
		var importingNodes []string
		for _, ms := range fsm.MasterSlaves {
			nodeId := ms.MasterId
			nodeSlots := fsm.Node2Slot[nodeId]
			if len(nodeSlots) > avgSlot+1 {
				exportingNodes = append(exportingNodes, nodeId)
			}
			if len(nodeSlots) < avgSlot-1 {
				importingNodes = append(importingNodes, nodeId)
			}
		}

		importIndex := 0
		exportIndex := 0
		var exportSlots []uint32
		for importIndex < len(importingNodes) && exportIndex < len(exportingNodes) {
			exportNode := exportingNodes[exportIndex]
			if len(exportSlots) == 0 {
				exportNodeSlots := fsm.Node2Slot[exportNode]
				exportCount := len(exportNodeSlots) - avgSlot
				exportSlots = exportNodeSlots[0:exportCount]
			}
			importNode := importingNodes[importIndex]
			importNodeCurrentIndex := fsm.Node2Slot[importNode]
			requirements := avgSlot - len(importNodeCurrentIndex)
			task := &raft.MigratingTask{
				ID:         utils.RandString(20),
				SrcNode:    exportNode,
				TargetNode: importNode,
			}
			if requirements <= len(exportSlots) {
				// exportSlots 可以提供足够的 slots, importingNode 处理完毕
				task.Slots = exportSlots[0:requirements]
				exportSlots = exportSlots[requirements:]
				importIndex++
			} else {
				// exportSlots 无法提供足够的 slots, exportingNode 处理完毕
				task.Slots = exportSlots
				exportSlots = nil
				exportIndex++
			}
			migratings = append(migratings, task)
		}
	})
	return migratings, nil
}

func (cluster *Cluster) waitCommitted(peer string, logIndex uint64) error {
	srcNodeConn, err := cluster.connections.BorrowPeerClient(peer)
	if err != nil {
		return err
	}
	defer cluster.connections.ReturnPeerClient(srcNodeConn)
	var peerIndex uint64
	for i := 0; i < 50; i++ {
		reply := srcNodeConn.Send(utils.ToCmdLine(getCommittedIndexCommand))
		switch reply := reply.(type) {
		case *protocol.IntReply:
			peerIndex = uint64(reply.Code)
			if peerIndex >= logIndex {
				return nil
			}
		case *protocol.StandardErrReply:
			logger.Infof("get committed index failed: %v", reply.Error())
		default:
			logger.Infof("get committed index unknown responseL %+v", reply.ToBytes())
		}
		time.Sleep(time.Millisecond * 100)
	}
	return errors.New("wait committed timeout")
}

// execMigrationChangeRoute should be exectued at leader
// it proposes EventFinishMigrate through raft, to change the route to the new node
// it returns until the proposal has been accepted by the majority  and two related nodes
// format: cluster.migration.changeroute taskid
func execMigrationChangeRoute(cluster *Cluster, c redis.Connection, cmdLine CmdLine) redis.Reply {
	if len(cmdLine) != 2 {
		return protocol.MakeArgNumErrReply(migrationChangeRouteCommand)
	}
	state := cluster.raftNode.State()
	if state != raft.Leader {
		// I am not leader, forward request to leader
		leaderConn, err := cluster.BorrowLeaderClient()
		if err != nil {
			return protocol.MakeErrReply(err.Error())
		}
		defer cluster.connections.ReturnPeerClient(leaderConn)
		return leaderConn.Send(cmdLine)
	}
	taskId := string(cmdLine[1])
	logger.Infof("change route for migration %s", taskId)
	task := cluster.raftNode.FSM.GetMigratingTask(taskId)
	if task == nil {
		return protocol.MakeErrReply("ERR task not found")
	}
	logger.Infof("change route for migration %s, got task info", taskId)
	// propose
	logIndex, err := cluster.raftNode.Propose(&raft.LogEntry{
		Event:         raft.EventFinishMigrate,
		MigratingTask: task,
	})
	if err != nil {
		return protocol.MakeErrReply("ERR " + err.Error())
	}
	logger.Infof("change route for migration %s, raft proposed", taskId)

	// confirm the 2 related node committed this log
	err = cluster.waitCommitted(task.SrcNode, logIndex)
	if err != nil {
		return protocol.MakeErrReply("ERR " + err.Error())
	}
	logger.Infof("change route for migration %s, confirm source node finished", taskId)

	err = cluster.waitCommitted(task.TargetNode, logIndex)
	if err != nil {
		return protocol.MakeErrReply("ERR " + err.Error())
	}
	logger.Infof("change route for migration %s, confirm target node finished", taskId)

	return protocol.MakeOkReply()
}
