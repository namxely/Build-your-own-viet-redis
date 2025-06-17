#!/bin/sh
set -e

# VietRedis Docker Entrypoint Script
# Handles initialization and configuration

echo "🇻🇳 Starting VietRedis Server Docker Container..."
echo "=================================================="

# Default values
VIETREDIS_CONFIG=${VIETREDIS_CONFIG:-"/etc/vietredis/vietredis.conf"}
VIETREDIS_DATA_DIR=${VIETREDIS_DATA_DIR:-"/data"}
VIETREDIS_LOG_LEVEL=${VIETREDIS_LOG_LEVEL:-"notice"}
VIETREDIS_LOG_FILE=${VIETREDIS_LOG_FILE:-"/var/log/vietredis/vietredis.log"}
VIETREDIS_MAX_MEMORY=${VIETREDIS_MAX_MEMORY:-""}
VIETREDIS_PORT=${VIETREDIS_PORT:-"6379"}

# Function to print messages
log_info() {
    echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_warn() {
    echo "[WARN] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_error() {
    echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Check if running as root
if [ "$(id -u)" = '0' ]; then
    log_info "Running as root, will switch to vietredis user"
    
    # Ensure proper ownership
    chown -R vietredis:vietredis "$VIETREDIS_DATA_DIR" /var/log/vietredis
    
    # Switch to vietredis user using su-exec
    exec su-exec vietredis "$0" "$@"
fi

# Ensure directories exist
mkdir -p "$VIETREDIS_DATA_DIR"
mkdir -p "$(dirname "$VIETREDIS_LOG_FILE")"

# Configuration file handling
if [ ! -f "$VIETREDIS_CONFIG" ]; then
    log_warn "Config file not found at $VIETREDIS_CONFIG, creating default config"
    
    cat > "$VIETREDIS_CONFIG" << EOF
# VietRedis Server Configuration (Auto-generated)
bind 0.0.0.0
port $VIETREDIS_PORT
timeout 0
databases 16
save 900 1
save 300 10
save 60 10000
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir $VIETREDIS_DATA_DIR
loglevel $VIETREDIS_LOG_LEVEL
logfile $VIETREDIS_LOG_FILE
appendonly yes
appendfilename vietredis.aof
appendfsync everysec
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
maxclients 10000
tcp-keepalive 300
tcp-backlog 511
EOF

    if [ -n "$VIETREDIS_MAX_MEMORY" ]; then
        echo "maxmemory $VIETREDIS_MAX_MEMORY" >> "$VIETREDIS_CONFIG"
        echo "maxmemory-policy allkeys-lru" >> "$VIETREDIS_CONFIG"
    fi
    
    if [ -n "$VIETREDIS_REPLICA_OF" ]; then
        echo "replicaof $VIETREDIS_REPLICA_OF" >> "$VIETREDIS_CONFIG"
        log_info "Configured as replica of $VIETREDIS_REPLICA_OF"
    fi
    
    if [ -n "$VIETREDIS_REQUIRE_PASS" ]; then
        echo "requirepass $VIETREDIS_REQUIRE_PASS" >> "$VIETREDIS_CONFIG"
        log_info "Password authentication enabled"
    fi
    
    if [ -n "$VIETREDIS_MASTER_AUTH" ]; then
        echo "masterauth $VIETREDIS_MASTER_AUTH" >> "$VIETREDIS_CONFIG"
        log_info "Master authentication configured"
    fi
fi

# Handle cluster mode
if [ "$VIETREDIS_CLUSTER_ENABLED" = "yes" ]; then
    log_info "Cluster mode enabled"
    echo "cluster-enabled yes" >> "$VIETREDIS_CONFIG"
    echo "cluster-config-file nodes.conf" >> "$VIETREDIS_CONFIG"
    echo "cluster-node-timeout 15000" >> "$VIETREDIS_CONFIG"
    
    if [ -n "$VIETREDIS_CLUSTER_ANNOUNCE_IP" ]; then
        echo "cluster-announce-ip $VIETREDIS_CLUSTER_ANNOUNCE_IP" >> "$VIETREDIS_CONFIG"
    fi
    
    if [ -n "$VIETREDIS_CLUSTER_ANNOUNCE_PORT" ]; then
        echo "cluster-announce-port $VIETREDIS_CLUSTER_ANNOUNCE_PORT" >> "$VIETREDIS_CONFIG"
    fi
fi

# Validate configuration
log_info "Validating configuration..."
if ! vietredis "$VIETREDIS_CONFIG" --test-config 2>/dev/null; then
    log_warn "Configuration validation failed, but continuing anyway"
fi

# Display startup information
log_info "VietRedis Server Configuration:"
log_info "  Config File: $VIETREDIS_CONFIG"
log_info "  Data Directory: $VIETREDIS_DATA_DIR"
log_info "  Log File: $VIETREDIS_LOG_FILE"
log_info "  Log Level: $VIETREDIS_LOG_LEVEL"
log_info "  Port: $VIETREDIS_PORT"
if [ -n "$VIETREDIS_MAX_MEMORY" ]; then
    log_info "  Max Memory: $VIETREDIS_MAX_MEMORY"
fi
if [ "$VIETREDIS_CLUSTER_ENABLED" = "yes" ]; then
    log_info "  Cluster Mode: Enabled"
fi

# Handle special cases
case "$1" in
    vietredis)
        log_info "Starting VietRedis Server..."
        exec "$@"
        ;;
    redis-cli)
        log_info "Starting Redis CLI..."
        exec "$@"
        ;;
    redis-benchmark)
        log_info "Starting Redis Benchmark..."
        exec "$@"
        ;;
    bash|sh)
        log_info "Starting shell..."
        exec "$@"
        ;;
    *)
        if [ -x "$1" ]; then
            log_info "Executing custom command: $*"
            exec "$@"
        else
            log_info "Starting VietRedis Server with config: $VIETREDIS_CONFIG"
            exec vietredis "$VIETREDIS_CONFIG"
        fi
        ;;
esac
