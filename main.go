package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

func main() {
	// Connect to Redis server
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Example of HMSET
	_, err = conn.Do("HMSET", "album:1",
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
