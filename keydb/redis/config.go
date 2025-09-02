package redis

import (
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
Config - if MasterName is not empty - FailOverClient will be created
if len Addresses more than 1 - ClusterClient will be created, except MasterName is set.
Else common redis.Client will be created.
*/
type Config struct {
	Addresses []string
	Password  string
	Database  int

	DialTimeout      time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	MaxConnectionAge time.Duration
	PoolTimeout      time.Duration
	IdleTimeout      time.Duration

	MaxRetries         int
	PoolSize           int
	MinIdleConnections int
	ReadOnly           bool
	MasterName         string
	SentinelPassword   string

	RouteByLatency bool
	RouteByRandom  bool
	// TLS
	TLS *tls.Config
}

func toUniversalRedisConfig(cfg Config) *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:            cfg.Addresses,
		DB:               cfg.Database,
		Password:         cfg.Password,
		MaxRetries:       cfg.MaxRetries,
		DialTimeout:      cfg.DialTimeout,
		ReadTimeout:      cfg.ReadTimeout,
		WriteTimeout:     cfg.WriteTimeout,
		PoolSize:         cfg.PoolSize,
		MinIdleConns:     cfg.MinIdleConnections,
		ConnMaxLifetime:  cfg.MaxConnectionAge,
		PoolTimeout:      cfg.PoolTimeout,
		ConnMaxIdleTime:  cfg.IdleTimeout,
		TLSConfig:        cfg.TLS,
		ReadOnly:         cfg.ReadOnly,
		RouteByLatency:   cfg.RouteByLatency,
		RouteRandomly:    cfg.RouteByRandom,
		MasterName:       cfg.MasterName,
		SentinelPassword: cfg.SentinelPassword,
	}
}

func toSentinelRedisConfig(cfg Config) *redis.Options {
	return &redis.Options{
		Addr:            cfg.Addresses[0],
		DB:              cfg.Database,
		Password:        cfg.Password,
		MaxRetries:      cfg.MaxRetries,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConnections,
		ConnMaxLifetime: cfg.MaxConnectionAge,
		PoolTimeout:     cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.IdleTimeout,
		TLSConfig:       cfg.TLS,
	}
}
