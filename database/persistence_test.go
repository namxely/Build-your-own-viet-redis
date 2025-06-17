package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/namxely/Build-your-own-viet-redis/aof"
	"github.com/namxely/Build-your-own-viet-redis/config"
	"github.com/namxely/Build-your-own-viet-redis/lib/utils"
	"github.com/namxely/Build-your-own-viet-redis/redis/connection"
	"github.com/namxely/Build-your-own-viet-redis/redis/protocol/asserts"
)

func TestLoadRDB(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(b))
	os.Chdir(projectRoot)
	config.Properties = &config.ServerProperties{
		AppendOnly:  false,
		RDBFilename: "test.rdb", // set working directory to project root
	}
	conn := connection.NewFakeConn()
	rdbDB := NewStandaloneServer()
	result := rdbDB.Exec(conn, utils.ToCmdLine("Get", "str"))
	asserts.AssertBulkReply(t, result, "str")
	result = rdbDB.Exec(conn, utils.ToCmdLine("LRange", "list", "0", "-1"))
	asserts.AssertMultiBulkReply(t, result, []string{"1", "2", "3", "4"})
	result = rdbDB.Exec(conn, utils.ToCmdLine("HGetAll", "hash"))
	asserts.AssertMultiBulkReply(t, result, []string{"1", "1"})
	result = rdbDB.Exec(conn, utils.ToCmdLine("ZRange", "zset", "0", "-1", "WITHSCORES"))
	asserts.AssertMultiBulkReply(t, result, []string{"0", "1", "1", "1"})
	result = rdbDB.Exec(conn, utils.ToCmdLine("SCard", "set"))
	asserts.AssertIntReply(t, result, 1)

	// test no rdb file
	config.Properties = &config.ServerProperties{
		AppendOnly:  false,
		RDBFilename:  "noexists.rdb", 
	}
	rdbDB = NewStandaloneServer()
	result = rdbDB.Exec(conn, utils.ToCmdLine("Get", "str"))
	asserts.AssertNullBulk(t, result)
}

func TestServerFsyncAlways(t *testing.T) {
	aofFile, err := ioutil.TempFile("", "*.aof")
	if err != nil {
		t.Error(err)
		return
	}
	config.Properties.AppendOnly = true
	config.Properties.AppendFilename = aofFile.Name()
	config.Properties.AppendFsync = aof.FsyncAlways
	server := NewStandaloneServer()
	conn := connection.NewFakeConn()
	server.Exec(conn, utils.ToCmdLine("del", "1"))
	ret := server.Exec(conn, utils.ToCmdLine("incr", "1"))
	asserts.AssertNotError(t, ret)
	reader := NewStandaloneServer()
	ret = reader.Exec(conn, utils.ToCmdLine("get", "1"))
	asserts.AssertBulkReply(t, ret, "1")
}

func TestServerFsyncEverySec(t *testing.T) {
	aofFile, err := ioutil.TempFile("", "*.aof")
	if err != nil {
		t.Error(err)
		return
	}
	config.Properties.AppendOnly = true
	config.Properties.AppendFilename = aofFile.Name()
	config.Properties.AppendFsync = aof.FsyncEverySec
	server := NewStandaloneServer()
	conn := connection.NewFakeConn()
	server.Exec(conn, utils.ToCmdLine("del", "1"))
	ret := server.Exec(conn, utils.ToCmdLine("incr", "1"))
	asserts.AssertNotError(t, ret)
	time.Sleep(1500 * time.Millisecond)
	reader := NewStandaloneServer()
	ret = reader.Exec(conn, utils.ToCmdLine("get", "1"))
	asserts.AssertBulkReply(t, ret, "1")
}
