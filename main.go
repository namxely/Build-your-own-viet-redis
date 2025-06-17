package main

/*
VietRedis Server - Enterprise Redis Implementation
=================================================

🇻🇳 Proudly Made in Vietnam by Namxely
📧 Contact: dev.namxely@gmail.com
🌐 Website: https://namxely.github.io/Build-your-own-viet-redis
📱 Telegram: @NamxelyDev
⭐ GitHub: https://github.com/namxely/Build-your-own-viet-redis

Copyright (c) 2025 Namxely Development
Licensed under MIT License

This high-performance Redis server implementation is designed
specifically for Vietnamese market with:
- Vietnamese language support
- Asia-Pacific optimized networking
- Vietnamese developer community focus
- Enterprise-grade reliability

Version: 2.0.0 "Saigon Edition"
Codename: "Phoenix Rising"
Release Date: June 17, 2025
Developer: Namxely (@namxely)
*/

import (
	"fmt"
	"os"
	"strings"

	"github.com/namxely/Build-your-own-viet-redis/cluster"
	"github.com/namxely/Build-your-own-viet-redis/config"
	"github.com/namxely/Build-your-own-viet-redis/database"
	idatabase "github.com/namxely/Build-your-own-viet-redis/interface/database"
	"github.com/namxely/Build-your-own-viet-redis/lib/logger"
	"github.com/namxely/Build-your-own-viet-redis/lib/utils"
	"github.com/namxely/Build-your-own-viet-redis/redis/server/gnet"
	stdserver "github.com/namxely/Build-your-own-viet-redis/redis/server/std"
)

// Build information - set during compile time
var (
	version   = "2.0.0-saigon"
	buildTime = "2025-06-17"
	gitCommit = "dev"
	buildBy   = "Namxely Dev Team"
)

// Server metadata
const (
	ServerName        = "VietRedis Server"
	ServerVersion     = "2.0.0"
	ServerCodename    = "Phoenix Rising"
	ServerEdition     = "Saigon Edition"
	ServerAuthor      = "Namxely (@namxely)"
	ServerWebsite     = "https://namxely.github.io/Build-your-own-viet-redis"
	ServerDescription = "🇻🇳 High-Performance Redis Implementation by Namxely"
)

var banner = `
╔══════════════════════════════════════════════╗
║        🚀 VietRedis Server v2.0 🚀           ║
║                                              ║
║   High-Performance Redis Implementation      ║
║   Developed by Namxely (@namxely)           ║
║   Built with ❤️  in Go Language             ║
║                                              ║
║   ⚡ Lightning Fast • 🔒 Secure • 📈 Scalable ║
║   🇻🇳 Made in Vietnam with Pride             ║
╚══════════════════════════════════════════════╝
`

var defaultProperties = &config.ServerProperties{
	Bind:           "0.0.0.0",
	Port:           6399,
	AppendOnly:     false,
	AppendFilename: "",
	MaxClients:     1000,
	RunID:          utils.RandString(40),
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

// printStartupInfo displays detailed server information
func printStartupInfo() {
	fmt.Println()
	fmt.Printf("🚀 Starting %s v%s (%s)\n", ServerName, ServerVersion, ServerCodename)
	fmt.Printf("📝 Edition: %s\n", ServerEdition)
	fmt.Printf("👨‍💻 Developed by: %s\n", ServerAuthor)
	fmt.Printf("🌐 Website: %s\n", ServerWebsite)
	fmt.Printf("📅 Build Date: %s\n", buildTime)
	fmt.Printf("🔧 Version: %s\n", version)
	fmt.Printf("👤 Built by: %s\n", buildBy)
	fmt.Println()
	fmt.Println("🔥 Features:")
	fmt.Println("   ✅ High-Performance Concurrent Engine")
	fmt.Println("   ✅ Full Redis Protocol Compatibility")
	fmt.Println("   ✅ Cluster Mode with Raft Consensus")
	fmt.Println("   ✅ AOF & RDB Persistence")
	fmt.Println("   ✅ Pub/Sub Messaging")
	fmt.Println("   ✅ Transactions & Pipelining")
	fmt.Println("   ✅ GEO Operations")
	fmt.Println("   ✅ Memory Optimization")
	fmt.Println()
	fmt.Println("🇻🇳 Proudly Made in Vietnam with ❤️")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()
}

func main() {
	// Display beautiful startup banner
	print(banner)
	printStartupInfo()

	// Setup logging with Vietnamese-friendly settings
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "vietredis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	configFilename := os.Getenv("CONFIG")
	if configFilename == "" {
		if fileExists("redis.conf") {
			config.SetupConfig("redis.conf")
		} else {
			config.Properties = defaultProperties
		}
	} else {
		config.SetupConfig(configFilename)
	}
	listenAddr := fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port)

	var err error
	if config.Properties.UseGnet {
		var db idatabase.DB
		if config.Properties.ClusterEnable {
			db = cluster.MakeCluster()
		} else {
			db = database.NewStandaloneServer()
		}
		server := gnet.NewGnetServer(db)
		err = server.Run(listenAddr)
	} else {
		handler := stdserver.MakeHandler()
		err = stdserver.Serve(listenAddr, handler)
	}
	if err != nil {
		logger.Errorf("start server failed: %v", err)
	}
}
