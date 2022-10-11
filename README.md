# Redsync

[![Go Reference](https://pkg.go.dev/badge/github.com/go-redsync/redsync/v4.svg)](https://pkg.go.dev/github.com/go-redsync/redsync/v4) [![Build Status](https://travis-ci.org/go-redsync/redsync.svg?branch=master)](https://travis-ci.org/go-redsync/redsync) 

Redlock provides a Redis-based distributed mutual exclusion lock implementation for Go as described in [this post](http://redis.io/topics/distlock). A reference library (by [antirez](https://github.com/antirez)) for Ruby is available at [github.com/antirez/redlock-rb](https://github.com/antirez/redlock-rb). </br>
Why i created this packet? </br>
This packet base from https://github.com/go-redis/redis. During the period of use, I needed access with quite high RQS to the lock, and the default config of https://github.com/go-redis/redis was unresponsive and unfriendly. Packet https://github.com/go-redis/redis has not updated for a long time, there are a few other packets that meet the requirements but are not as strong as the community test suite same https://github.com/go-redis/redis. I need a simple custom Mutex, easy to use for my project. I was customs Mutex for my project. </br>
## Installation

Install  Redsync using the go get command:

    $ go get github.com/Nghiait123456/redlock

Two driver    implementations will be installed; however, only the one used will be included in your project.

 * [Redigo](https://github.com/gomodule/redigo)
 * [Go-redis](https://github.com/go-redis/redis)

See the [examples](examples) folder for usage of each driver.

## Documentation

- [Reference](https://godoc.org/github.com/Nghiait123456/redlock)

## Usage

Redis-cluster example:

```go
package main

import (
	"context"
	"fmt"
	"github.com/Nghiait123456/redlock"
	"github.com/Nghiait123456/redlock/redis/goredis/v8"
	goredislib "github.com/go-redis/redis/v8"
	"time"
)

func main() {
	client := goredislib.NewClusterClient(&goredislib.ClusterOptions{
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "bitnami",
	})

	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	mutex := rs.NewMutex("test-redsync")
	ctx := context.Background()

	fmt.Println("start lock")
	if err := mutex.LockContext(ctx); err != nil {
		fmt.Println("lock fail")
		panic(err)
	}
	fmt.Println("start lock success")

	fmt.Println("start race condition lock 1st")
	go func() {
		fmt.Println("start race conditions lock 1st")
		if err := mutex.LockContext(ctx); err != nil {
			fmt.Printf("race conditions fail 1st, err: %v \n", err.Error())
		}
		fmt.Println("race conditions lock success 1st")
	}()

	time.Sleep(10 * time.Second)

	fmt.Println("start end lock")
	if _, err := mutex.UnlockContext(ctx); err != nil {
		fmt.Printf("race conditions fail 1st, err: %v \n", err.Error())
		panic(err)
	}

	fmt.Println("start race condition lock 2st")
	go func() {
		fmt.Println("start race conditions lock 2st")
		if err := mutex.LockContext(ctx); err != nil {
			fmt.Println("race conditions fail 2st")
			panic(err)
		}
		fmt.Println("race conditions lock success 2st")
	}()

	time.Sleep(1 * time.Second)

	fmt.Println("end lock success")
}
```


All-example in  [Link](https://github.com/Nghiait123456/RedLock/tree/master/examples)