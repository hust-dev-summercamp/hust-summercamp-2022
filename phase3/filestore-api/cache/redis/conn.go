package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

// newRedisPool : 创建redis连接池
func newRedisPool(redisHost string, redisPass string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 1. 打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// 2. 访问认证
			if _, err = c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func InitRedis(redisHost string, redisPass string) {
	fmt.Println(redisHost, redisPass)
	pool = newRedisPool(redisHost, redisPass)
}

//
//func init(redisHost1 string, redisPass1 string) {
//	pool = newRedisPool(os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))
//}

func RedisPool() *redis.Pool {
	return pool
}
