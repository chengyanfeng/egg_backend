package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//redis 添加
func ReAdd(key, value string, time int) {

	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()

	_, err = c.Do("SET", key, value, "EX", time)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

}

//redis 获取
func ReGet(key string) interface{} {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return false
	}
	defer c.Close()

	username, err := redis.String(c.Do("GET", key))
	if err != nil {
		fmt.Println("redis get failed:", err)
		return ""
	} else {
		fmt.Printf("Get mykey: %v \n", username)
		return username

	}
}

//检查key 是否存在
func ReIsEx(key string) bool {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return false
	}
	defer c.Close()

	is_key_exit, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		fmt.Println("error:", err)
		return false
	} else {
		fmt.Printf("exists or not: %v \n", is_key_exit)
		return true
	}
}

//删除
func ReDel(key string) bool {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return false
	}
	defer c.Close()
	_, err = c.Do("DEL", key)
	if err != nil {
		fmt.Println("redis delelte failed:", err)
		return false
	}
	return true
}

//设置某个key的过期时间
func ReExpr(key string, time int) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()

	// 设置过期时间为24小时
	n, _ := c.Do("EXPIRE", key, 24*time)
	if n == int64(1) {
		fmt.Println("success")
	}

}
