package gocache

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/go-redis/redis"
)

// Config represents the cache configuration to
// be used depending of the specified backend.
// Only one backend should be specified per
// cache meaning that the backend config
// should not be nil
type Config struct {
	Redis    *RedisConfig
	Memcache *MemcacheConfig
	Map      *MapConfig
}

// RedisConfig represents the configuration
// for a cache with a redis backend
type RedisConfig struct {
	// The prefix to be appended to every cache entry
	Prefix string
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string
	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer func() (net.Conn, error)
	// Hook that is called when new connection is established.
	OnConnect func(*redis.Conn) error
	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string
	// Database to be selected after connecting to the server.
	DB int
	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration
	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration
	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config
}

// MemcacheConfig represents the configuration
// for a cache with a Memcache backend
type MemcacheConfig struct {
	// The prefix to be appended to every cache entry
	Prefix string
	// Timeout specifies the socket read/write timeout.
	// If zero, DefaultTimeout is used.
	Timeout time.Duration
	// MaxIdleConns specifies the maximum number of idle connections that will
	// be maintained per address. If less than one, DefaultMaxIdleConns will be
	// used.
	//
	// Consider your expected traffic rates and latency carefully. This should
	// be set to a number higher than your peak parallel requests.
	MaxIdleConns int
	// Server list to be used by the client. Each server is weighted the same
	Servers []string
}

// MapConfig represents the configuration
// for a cache with a map backend
type MapConfig struct {
	// The prefix to be appended to every cache entry
	Prefix string
}
