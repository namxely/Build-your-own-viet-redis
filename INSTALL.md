# Hướng dẫn cài đặt và chạy VietRedis Server

## 🇻🇳 Giới thiệu

VietRedis Server là Redis implementation được phát triển bởi **Namxely** - một Vietnamese developer, tối ưu hóa cho thị trường Việt Nam và châu Á.

## Yêu cầu hệ thống

- **Go 1.18+**: Tải từ [https://golang.org/dl/](https://golang.org/dl/)
- **Git**: Để clone repository
- **Redis CLI**: (Tùy chọn) Để test kết nối

## Cài đặt Go (Windows)

1. Tải Go installer từ [https://golang.org/dl/](https://golang.org/dl/)
2. Chạy file `.msi` và làm theo hướng dẫn
3. Mở Command Prompt hoặc PowerShell mới
4. Kiểm tra cài đặt: `go version`

## Cài đặt và chạy project

### Cách 1: Clone repository

```bash
# Clone repository
git clone https://github.com/namxely/viet-redis-server.git
cd viet-redis-server

# Tải dependencies
go mod tidy

# Build
go build -o vietredis.exe

# Chạy
./vietredis.exe
```

### Cách 2: Sử dụng build script

**Windows PowerShell:**
```powershell
# Cấp quyền execute (chỉ cần 1 lần)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Chạy build script
.\build.ps1
```

**Linux/macOS:**
```bash
chmod +x build.sh
./build.sh
```

### Cách 3: Chạy trực tiếp

```bash
go run main.go
```

## Cấu hình

### File cấu hình mặc định (redis.conf)

```
# Network
bind 0.0.0.0
port 6399

# General
databases 16
timeout 0

# Persistence
save 900 1
save 300 10
save 60 10000

# AOF
appendonly no
appendfilename "appendonly.aof"

# Logging
loglevel notice

# Cluster
cluster-enabled no
```

### Chạy với cấu hình tùy chỉnh

```bash
./my-redis-server.exe redis.conf
```

## Test kết nối

### Sử dụng Redis CLI

```bash
# Cài đặt Redis CLI (Windows - Chocolatey)
choco install redis-64

# Hoặc tải từ https://github.com/microsoftarchive/redis/releases

# Kết nối
redis-cli -p 6399

# Test commands
127.0.0.1:6399> SET hello world
OK
127.0.0.1:6399> GET hello
"world"
```

### Sử dụng Telnet

```bash
telnet localhost 6399

# Test commands
SET hello world
GET hello
```

### Sử dụng curl (REST API - nếu enabled)

```bash
curl -X POST "http://localhost:6399/api/v1/set" -d '{"key":"hello","value":"world"}'
curl "http://localhost:6399/api/v1/get?key=hello"
```

## Chế độ Cluster

### Cấu hình cluster

**node1.conf:**
```
port 7000
cluster-enabled yes
cluster-config-file nodes-7000.conf
cluster-node-timeout 5000
```

**node2.conf:**
```
port 7001
cluster-enabled yes
cluster-config-file nodes-7001.conf
cluster-node-timeout 5000
```

**node3.conf:**
```
port 7002
cluster-enabled yes
cluster-config-file nodes-7002.conf
cluster-node-timeout 5000
```

### Chạy cluster

```bash
# Terminal 1
./my-redis-server.exe node1.conf

# Terminal 2
./my-redis-server.exe node2.conf

# Terminal 3
./my-redis-server.exe node3.conf
```

## Troubleshooting

### Lỗi thường gặp

**1. "go: command not found"**
- Cài đặt Go từ [golang.org](https://golang.org/dl/)
- Thêm Go vào PATH environment variable

**2. "Port already in use"**
- Thay đổi port trong config file
- Hoặc kill process đang sử dụng port: `netstat -ano | findstr :6399`

**3. "Permission denied"**
- Windows: Chạy với quyền Administrator
- Linux/macOS: `sudo chmod +x my-redis-server`

**4. Import cycle errors**
- Chạy `go mod tidy` để clean dependencies
- Kiểm tra circular imports trong code

### Debug

**Enable debug logging:**
```bash
./my-redis-server.exe --loglevel debug
```

**Monitor connections:**
```bash
# Redis CLI
redis-cli -p 6399 MONITOR
```

**Check server info:**
```bash
redis-cli -p 6399 INFO
```

## Performance Tuning

### System limits

**Linux:**
```bash
# Increase file descriptor limit
ulimit -n 65535

# Increase TCP backlog
echo 511 > /proc/sys/net/core/somaxconn
```

**Windows:**
```powershell
# Increase TCP settings in registry
# HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters
```

### Application tuning

```
# redis.conf
tcp-backlog 511
tcp-keepalive 300
timeout 0
maxclients 10000
```

## Monitoring

### Built-in commands

```
INFO        # Server info
CLIENT LIST # Connected clients
SLOWLOG GET # Slow queries
MEMORY USAGE key # Memory usage
```

### External tools

- **redis-stat**: Real-time Redis monitoring
- **RedisInsight**: GUI tool
- **Grafana + Prometheus**: Metrics dashboard

## Backup và Recovery

### RDB Snapshots

```bash
# Manual snapshot
redis-cli -p 6399 BGSAVE

# Automatic snapshots (redis.conf)
save 900 1      # Save if at least 1 key changed in 900 seconds
save 300 10     # Save if at least 10 keys changed in 300 seconds
save 60 10000   # Save if at least 10000 keys changed in 60 seconds
```

### AOF (Append Only File)

```bash
# Enable AOF (redis.conf)
appendonly yes
appendfilename "appendonly.aof"

# Manual rewrite
redis-cli -p 6399 BGREWRITEAOF
```

## Development

### Code structure

```
main.go              # Entry point
├── config/          # Configuration
├── database/        # Core database logic
├── redis/           # Redis protocol
├── cluster/         # Cluster implementation
├── datastruct/      # Data structures
├── lib/             # Utilities
└── tcp/             # Network layer
```

### Adding new commands

1. Define command in `database/router.go`
2. Implement logic in appropriate file
3. Add tests
4. Update documentation

### Running tests

```bash
# All tests
go test ./...

# Specific package
go test ./database

# With coverage
go test -cover ./...

# Race detection
go test -race ./...
```

---

**Liên hệ hỗ trợ:**
- Author: Namxely (@namxely)
- Email: dev.namxely@gmail.com
- GitHub Issues: [Create issue](https://github.com/namxely/viet-redis-server/issues)
- Telegram: @NamxelyDev
