package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

var redisPool *redis.Pool

func initRedisPool() {
	redisPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Fatal(err)
			}
			return conn, err
		},
	}
}

func hmset(key string, args ...interface{}) error {
	log.Printf("DEBUG: HMSET %s %v", key, args)

	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(args)...)
	if err != nil {
		log.Printf("ERROR: fail HMSET %s %s %s", key, args, err.Error())
		return err
	}

	return nil
}

func main() {
	// Connect to Redis server
	initRedisPool()
	conn := redisPool.Get()

	defer conn.Close()

	// Example of HMSET
	err := hmset("album:1",
		"title", "Back in Black",
		"artist", "AC/DC",
		"year-released", 1980,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Album added")

	// Example of HMGET
	title, err := redis.String(conn.Do("HGET", "album:1", "title"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Title: %s\n", title)

	// Example of HMGETALL and Values/ScanStruct
	type Album struct {
		Title        string `redis:"title"`
		Artist       string `redis:"artist"`
		YearReleased int    `redis:"year-released"`
	}

	values, err := redis.Values(conn.Do("HGETALL", "album:1"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", values)

	var album Album
	err = redis.ScanStruct(values, &album)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", album)
}
