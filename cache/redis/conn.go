package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
)

// 创建连接池 链接redis
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			//1打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return nil, err
			}
			// 访问认证 输入密码
			return c, nil
		},
		// 检测连接是否失效
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil //如果小于一分钟  检测
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}
func init() {
	pool = newRedisPool()
}
func RedisPool() *redis.Pool {
	return pool
}
