package database

import (
	"github.com/namxely/Build-your-own-viet-redis/config"
	"github.com/namxely/Build-your-own-viet-redis/lib/utils"
	"github.com/namxely/Build-your-own-viet-redis/redis/connection"
	"github.com/namxely/Build-your-own-viet-redis/redis/protocol/asserts"
	"math/rand"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	c := connection.NewFakeConn()
	actual := Ping(c, utils.ToCmdLine())
	asserts.AssertStatusReply(t, actual, "PONG")
	val := utils.RandString(5)
	actual = Ping(c, utils.ToCmdLine(val))
	asserts.AssertStatusReply(t, actual, val)
	actual = Ping(c, utils.ToCmdLine(val, val))
	asserts.AssertErrReply(t, actual, "ERR wrong number of arguments for 'ping' command")
}

func TestAuth(t *testing.T) {
	passwd := utils.RandString(10)
	c := connection.NewFakeConn()
	ret := testServer.Exec(c, utils.ToCmdLine("AUTH"))
	asserts.AssertErrReply(t, ret, "ERR wrong number of arguments for 'auth' command")
	ret = testServer.Exec(c, utils.ToCmdLine("AUTH", passwd))
	asserts.AssertErrReply(t, ret, "ERR Client sent AUTH, but no password is set")

	config.Properties.RequirePass = passwd
	defer func() {
		config.Properties.RequirePass = ""
	}()
	ret = testServer.Exec(c, utils.ToCmdLine("AUTH", passwd+"wrong"))
	asserts.AssertErrReply(t, ret, "ERR invalid password")
	ret = testServer.Exec(c, utils.ToCmdLine("GET", "A"))
	asserts.AssertErrReply(t, ret, "NOAUTH Authentication required")
	ret = testServer.Exec(c, utils.ToCmdLine("AUTH", passwd))
	asserts.AssertStatusReply(t, ret, "OK")

}

func TestInfo(t *testing.T) {
	c := connection.NewFakeConn()
	ret := testServer.Exec(c, utils.ToCmdLine("INFO"))
	asserts.AssertNotError(t, ret)
	ret = testServer.Exec(c, utils.ToCmdLine("INFO", "server"))
	asserts.AssertNotError(t, ret)
	ret = testServer.Exec(c, utils.ToCmdLine("INFO", "client"))
	asserts.AssertNotError(t, ret)
	ret = testServer.Exec(c, utils.ToCmdLine("INFO", "cluster"))
	asserts.AssertNotError(t, ret)
	ret = testServer.Exec(c, utils.ToCmdLine("iNFO", "SeRvEr"))
	asserts.AssertNotError(t, ret)
	ret = testServer.Exec(c, utils.ToCmdLine("INFO", "Keyspace"))
	asserts.AssertNotError(t, ret)
	ret = testServer.Exec(c, utils.ToCmdLine("iNFO", "abc", "bde"))
	asserts.AssertErrReply(t, ret, "ERR wrong number of arguments for 'info' command")
	ret = testServer.Exec(c, utils.ToCmdLine("INFO", "abc"))
	asserts.AssertErrReply(t, ret, "Invalid section for 'info' command")
}

func TestDbSize(t *testing.T) {
	c := connection.NewFakeConn()
	rand.NewSource(time.Now().UnixNano())
	randomNum := rand.Intn(10) + 1
	for i := 0; i < randomNum; i++ {
		key := utils.RandString(10)
		value := utils.RandString(10)
		testServer.Exec(c, utils.ToCmdLine("SET", key, value))
	}
	ret := testServer.Exec(c, utils.ToCmdLine("dbsize"))
	asserts.AssertIntReply(t, ret, randomNum)
}
