package commands

import (
	"github.com/namxely/Build-your-own-viet-redis/cluster/core"
	"github.com/namxely/Build-your-own-viet-redis/interface/redis"
	"github.com/namxely/Build-your-own-viet-redis/lib/utils"
	"github.com/namxely/Build-your-own-viet-redis/redis/protocol"
)

func init() {
	core.RegisterCmd("rename_", execRenameInLocal)
	core.RegisterCmd("renamenx_", execRenameNxInLocal)
	core.RegisterCmd("rename", execRename)
	core.RegisterCmd("renamenx", execRenameNx)
	core.RegisterPrepareFunc("RenameFrom", prepareRenameFrom)
	core.RegisterPrepareFunc("RenameNxTo", prepareRenameNxTo)

}

func execRenameInLocal(cluster *core.Cluster, c redis.Connection, cmdLine CmdLine) redis.Reply {
	cmdLine[0] = []byte("rename")
	return cluster.LocalExec(c, cmdLine)
}

func execRenameNxInLocal(cluster *core.Cluster, c redis.Connection, cmdLine CmdLine) redis.Reply {
	cmdLine[0] = []byte("renamenx")
	return cluster.LocalExec(c, cmdLine)
}

func execRename(cluster *core.Cluster, c redis.Connection, cmdLine CmdLine) redis.Reply {
	if len(cmdLine) != 3 {
		return protocol.MakeArgNumErrReply("rename")
	}
	src := string(cmdLine[1])
	target := string(cmdLine[2])
	srcSlot := cluster.GetSlot(src)
	srcNode := cluster.PickNode(srcSlot)
	targetSlot := cluster.GetSlot(target)
	targetNode := cluster.PickNode(targetSlot)
	if srcNode == targetNode {
		cmdLine[0] = []byte("rename")
		return cluster.Relay(srcNode, c, cmdLine)
	}
	routeMap := RouteMap{
		srcNode:    {src},
		targetNode: {target},
	}

	txID := utils.RandString(10)
	srcPrepareResp := cluster.Relay(srcNode, c, utils.ToCmdLine("Prepare", txID, "RenameFrom", src))
	if protocol.IsErrorReply(srcPrepareResp) {
		// rollback src node
		requestRollback(cluster, c, txID, map[string][]string{srcNode: {src}})
		return srcPrepareResp
	}
	srcPrepareResult, ok := srcPrepareResp.(*protocol.MultiBulkReply)
	if !ok || len(srcPrepareResult.Args) < 2 {
		requestRollback(cluster, c, txID, map[string][]string{srcNode: {src}})
		return protocol.MakeErrReply("ERR invalid prepare response")
	}
	// prepare rename to
	targetPrepareResp := cluster.Relay(targetNode, c, utils.ToCmdLine3("Prepare", []byte(txID),
		[]byte("RenameTo"), []byte(target), srcPrepareResult.Args[0], srcPrepareResult.Args[1]))
	if protocol.IsErrorReply(targetPrepareResp) {
		// rollback src node
		requestRollback(cluster, c, txID, routeMap)
		return targetPrepareResp
	}

	// commit
	commitCmd := utils.ToCmdLine("commit", txID)
	for node := range routeMap {
		reply := cluster.Relay(node, c, commitCmd)
		if err := protocol.Try2ErrorReply(reply); err != nil {
			requestRollback(cluster, c, txID, routeMap)
			return protocol.MakeErrReply("commit failed: " + err.Error())
		}
	}
	return protocol.MakeOkReply()
}

func execRenameNx(cluster *core.Cluster, c redis.Connection, cmdLine CmdLine) redis.Reply {
	if len(cmdLine) != 3 {
		return protocol.MakeArgNumErrReply("rename")
	}
	src := string(cmdLine[1])
	target := string(cmdLine[2])
	srcSlot := cluster.GetSlot(src)
	srcNode := cluster.PickNode(srcSlot)
	targetSlot := cluster.GetSlot(target)
	targetNode := cluster.PickNode(targetSlot)
	if srcNode == targetNode {
		cmdLine[0] = []byte("rename")
		return cluster.Relay(srcNode, c, cmdLine)
	}
	routeMap := RouteMap{
		srcNode:    {src},
		targetNode: {target},
	}

	txID := utils.RandString(10)
	srcPrepareResp := cluster.Relay(srcNode, c, utils.ToCmdLine("Prepare", txID, "RenameFrom", src))
	if protocol.IsErrorReply(srcPrepareResp) {
		// rollback src node
		requestRollback(cluster, c, txID, map[string][]string{srcNode: {src}})
		return srcPrepareResp
	}
	srcPrepareResult, ok := srcPrepareResp.(*protocol.MultiBulkReply)
	if !ok || len(srcPrepareResult.Args) < 2 {
		requestRollback(cluster, c, txID, map[string][]string{srcNode: {src}})
		return protocol.MakeErrReply("ERR invalid prepare response")
	}
	// prepare rename to
	targetPrepareResp := cluster.Relay(targetNode, c, utils.ToCmdLine3("Prepare", []byte(txID),
		[]byte("RenameNxTo"), []byte(target), srcPrepareResult.Args[0], srcPrepareResult.Args[1]))
	if res, ok := targetPrepareResp.(protocol.ErrorReply); ok {
		// rollback src node
		requestRollback(cluster, c, txID, routeMap)
		if res.Error() == keyExistsErr {
			return protocol.MakeIntReply(0)
		}
		return res
	}

	// commit
	commitCmd := utils.ToCmdLine("commit", txID)
	for node := range routeMap {
		reply := cluster.Relay(node, c, commitCmd)
		if err := protocol.Try2ErrorReply(reply); err != nil {
			requestRollback(cluster, c, txID, routeMap)
			return protocol.MakeErrReply("commit failed: " + err.Error())
		}
	}
	return protocol.MakeIntReply(1)
}

// prepareRenameFrom is prepare-function for RenameFrom
func prepareRenameFrom(cluster *core.Cluster, conn redis.Connection, cmdLine CmdLine) redis.Reply {
	if len(cmdLine) != 2 {
		return protocol.MakeArgNumErrReply("RenameFrom")
	}
	key := string(cmdLine[1])
	existResp := cluster.LocalExec(conn, utils.ToCmdLine("Exists", key))
	if protocol.IsErrorReply(existResp) {
		return existResp
	}
	existIntResp := existResp.(*protocol.IntReply)
	if existIntResp.Code == 0 {
		return protocol.MakeErrReply("ERR no such key")
	}
	return cluster.LocalExecWithinLock(conn, utils.ToCmdLine2("DumpKey", key))
}

const keyExistsErr = "key exists"

func prepareRenameNxTo(cluster *core.Cluster, conn redis.Connection, cmdLine CmdLine) redis.Reply {
	if len(cmdLine) != 4 {
		return protocol.MakeArgNumErrReply("RenameNxTo")
	}
	key := string(cmdLine[1])
	exists := cluster.LocalExists([]string{key})
	if len(exists) > 0 {
		return protocol.MakeErrReply(keyExistsErr)
	}
	return protocol.MakeOkReply()
}
