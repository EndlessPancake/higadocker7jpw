package main

import (
        "fmt"
        "github.com/gomodule/redigo/redis"
)

func main() {
        c, err := redis.Dial("tcp", "0.0.0.0:6379")
        if err != nil {
                fmt.Println(err)
                return
        }
        defer c.Close()
        fmt.Println(c.Do("hset", "my_hash", "key1", "value1"))
        fmt.Println(c.Do("hset", "my_hash", "key2", "value2"))
        fmt.Println(redis.String(c.Do("hget", "my_hash", "key1")))
        myMap, err := redis.StringMap(c.Do("hgetall", "my_hash"))
        if err != nil {
                fmt.Println(err)
                return
        }
        for k, v := range myMap {
                fmt.Printf("%s -> %s\n", k, v)
        }
}
