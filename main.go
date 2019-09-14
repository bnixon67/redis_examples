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
	_, err = conn.Do("HMSET", "album:2",
		"title", "Electric Ladyland",
		"artist", "Jimi Hendrix",
		"price", 4.95,
		"likes", 8)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Electric Ladyland added!")

	// Example of HMGET
	title, err := redis.String(conn.Do("HGET", "album:2", "title"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Title: %s\n", title)

	// Example of HMGETALL and Values/ScanStruct
	type Album struct {
		Title  string  `redis:"title"`
		Artist string  `redis:"artist"`
		Price  float64 `redis:"price"`
		Likes  int     `redis:"likes"`
	}

	values, err := redis.Values(conn.Do("HGETALL", "album:2"))
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
