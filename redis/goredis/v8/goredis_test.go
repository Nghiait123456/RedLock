package goredis

import "github.com/Nghiait123456/redlock/redis"

var _ redis.Conn = (*conn)(nil)

var _ redis.Pool = (*pool)(nil)
