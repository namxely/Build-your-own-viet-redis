// Package cluster provides a server side cluster which is transparent to client. You can connect to any node in the cluster to access all data in the cluster
package cluster

import (
	"os"
	"path"

	_ "github.com/namxely/Build-your-own-viet-redis/cluster/commands" // register commands
	"github.com/namxely/Build-your-own-viet-redis/cluster/core"
	"github.com/namxely/Build-your-own-viet-redis/cluster/raft"
	"github.com/namxely/Build-your-own-viet-redis/config"
	"github.com/namxely/Build-your-own-viet-redis/lib/logger"
)

type Cluster = core.Cluster

// MakeCluster creates and starts a node of cluster
func MakeCluster() *Cluster {
	raftPath := path.Join(config.Properties.Dir, "raft")
	err := os.MkdirAll(raftPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	cluster, err := core.NewCluster(&core.Config{
		RaftConfig: raft.RaftConfig{
			RedisAdvertiseAddr: config.Properties.AnnounceAddress(),
			RaftListenAddr:     config.Properties.RaftListenAddr,
			RaftAdvertiseAddr:  config.Properties.RaftAnnounceAddress(),
			Dir:                raftPath,
		},
		StartAsSeed: config.Properties.ClusterAsSeed,
		JoinAddress: config.Properties.ClusterSeed,
		Master:      config.Properties.MasterInCluster,
	})
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	return cluster
}
