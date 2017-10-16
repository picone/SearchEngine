package redigo

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 3 * time.Minute,
		Dial: func() (redis.Conn, error){
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}

func GetConnection() redis.Conn {
	return pool.Get()
}
