package std

/*
 * A tcp.Handler implements redis protocol
 */

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/namxely/Build-your-own-viet-redis/cluster"
	"github.com/namxely/Build-your-own-viet-redis/config"
	"github.com/namxely/Build-your-own-viet-redis/database"
	idatabase "github.com/namxely/Build-your-own-viet-redis/interface/database"
	"github.com/namxely/Build-your-own-viet-redis/lib/logger"
	"github.com/namxely/Build-your-own-viet-redis/lib/sync/atomic"
	"github.com/namxely/Build-your-own-viet-redis/redis/connection"
	"github.com/namxely/Build-your-own-viet-redis/redis/parser"
	"github.com/namxely/Build-your-own-viet-redis/redis/protocol"
	"github.com/namxely/Build-your-own-viet-redis/tcp"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

// Handler implements tcp.Handler and serves as a redis server
type Handler struct {
	activeConn sync.Map // *client -> placeholder
	db         idatabase.DB
	closing    atomic.Boolean // refusing new client and new request
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handler {
	var db idatabase.DB
	if config.Properties.ClusterEnable {
		db = cluster.MakeCluster()
	} else {
		db = database.NewStandaloneServer()
	}
	return &Handler{
		db: db,
	}
}

func Serve(addr string, handler *Handler) error {
	return tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: addr,
	}, handler)
}

func (h *Handler) closeClient(client *connection.Connection) {
	_ = client.Close()
	h.db.AfterClientClose(client)
	h.activeConn.Delete(client)
}


// Handle receives and executes redis commands
func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
		return
	}

	client := connection.NewConn(conn)
	h.activeConn.Store(client, struct{}{})

	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// connection closed
				h.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr())
				return
			}
			// protocol err
			errReply := protocol.MakeErrReply(payload.Err.Error())
			_, err := client.Write(errReply.ToBytes())
			if err != nil {
				h.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr())
				return
			}
			continue
		}
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := payload.Data.(*protocol.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk protocol")
			continue
		}
		result := h.db.Exec(client, r.Args)
		if result != nil {
			_, _ = client.Write(result.ToBytes())
		} else {
			_, _ = client.Write(unknownErrReplyBytes)
		}
	}
}

// Close stops handler
func (h *Handler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	// TODO: concurrent wait
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	h.db.Close()
	return nil
}
