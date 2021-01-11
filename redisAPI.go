package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Entry struct {
	Uid      string  `json:"UserID"`
	Name     string  `json:"Name"`
	Author   string  `json:"Author"`
	Score    float32 `json:"Score"`
	Url      string  `json:"URL"`
	Synopsis string  `json:"Synopsis"`
	Comments string  `json:"Comments"`
	Lastupdt string  `json:"LastUpdate"`
	Category string  `json:"category"`
}

type EntryKey struct {
	Uid    string `json:"UserID"`
	Name   string `json:"Name"`
	Author string `json:"Author"`
}

var ctx = context.Background()

var rdb1 *redis.Client //note
var rdb2 *redis.Client //trash

func InitRedis() {
	rdb1 = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,
	})

	_, err := rdb1.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	rdb2 = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       2,
	})

	_, err = rdb2.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func InsertEntry(entrystr string) {
	// error check
	if rdb1 == nil || rdb2 == nil {
		panic("Redis not connected")
	}
	et := &Entry{}
	err := json.Unmarshal([]byte(entrystr), et)
	if err != nil {
		fmt.Println(err)
		return
	}

	//get key
	ek := &EntryKey{et.Uid, et.Name, et.Author}
	ekey, _ := json.Marshal(ek)
	key := string(ekey)

	// No duplicate
	_, err = rdb1.Get(ctx, key).Result()
	if err != redis.Nil {
		fmt.Println("InsertEntry already exists Or other error occurs")
		return
	}

	// Insert entry
	err = rdb1.Set(ctx, key, entrystr, 0).Err()
	if err != nil {
		fmt.Println("InsertEntry HSet error")
		return
	}
}

func RemoveEntry(entrystr string) {
	// error check
	if rdb1 == nil || rdb2 == nil {
		panic("Redis not connected")
	}
	et := &Entry{}
	err := json.Unmarshal([]byte(entrystr), et)
	if err != nil {
		fmt.Println(err)
		return
	}

	//get key
	ek := &EntryKey{et.Uid, et.Name, et.Author}
	ekey, _ := json.Marshal(ek)
	key := string(ekey)

	// Avoid new entry
	_, err = rdb1.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		fmt.Println("RemoveEntry never exists")
		return
	}

	// Remove entry
	err = rdb1.Del(ctx, key).Err()
	if err != nil {
		fmt.Println("Remove HDel error")
		return
	}

	// Move to Trash
	err = rdb2.Set(ctx, key, entrystr, 14*24*time.Hour).Err()
	if err != nil {
		fmt.Println("Remove to trash error")
		return
	}
}

func RecoverEntry(entrystr string) {
	// error check
	if rdb1 == nil || rdb2 == nil {
		panic("Redis not connected")
	}
	et := &Entry{}
	err := json.Unmarshal([]byte(entrystr), et)
	if err != nil {
		fmt.Println(err)
		return
	}

	//get key
	ek := &EntryKey{et.Uid, et.Name, et.Author}
	ekey, _ := json.Marshal(ek)
	key := string(ekey)

	// Avoid new entry
	_, err = rdb2.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		fmt.Println("RecoverEntry never exists")
		return
	}

	// Remove entry
	err = rdb2.Del(ctx, key).Err()
	if err != nil {
		fmt.Println("Recover HDel error")
		return
	}

	// Move to note
	err = rdb1.Set(ctx, key, entrystr, 0).Err()
	if err != nil {
		fmt.Println("Recover to Note error")
		return
	}
}

func GetAll() string {
	// error check
	if rdb1 == nil || rdb2 == nil {
		panic("Redis not connected")
	}
	var finalstr = ""
	//Get all notes
	iter := rdb1.Scan(ctx, 0, "", 0).Iterator()
	if err := iter.Err(); err != nil {
		panic(err)
	}
	for iter.Next(ctx) {
		key := iter.Val()
		etstr, _ := rdb1.Get(ctx, key).Result()
		finalstr += "Note:" + etstr + "\n"
	}

	//Get all trashes
	iter2 := rdb2.Scan(ctx, 0, "", 0).Iterator()
	if err := iter2.Err(); err != nil {
		panic(err)
	}
	for iter2.Next(ctx) {
		key := iter2.Val()
		etstr, _ := rdb2.Get(ctx, key).Result()
		finalstr += "Trash:" + etstr + "\n"
	}
	return finalstr
}
