package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

// 初始化连接池
var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = "137908"
)

// 创建连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 打开连接
			conn, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}
			// 设置密码
			_, err = conn.Do("AUTH", redisPass)
			if err != nil {
				fmt.Println("连接失败")
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		// 长连接，每个一分钟查看是否还链接成功
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			// 测额是下ping值
			_, err := conn.Do("PING")
			return err
		},
	}
}

// 初始化就创建连接池
func init() {
	pool = newRedisPool()
	// 初始化redis时查所有的keys
	data, err := pool.Get().Do("KEYS", "*")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

// 向外暴露连接池的接口
func GetRedisPool() *redis.Pool {
	return pool
}
