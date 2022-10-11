package main

import (
	"context"
	"fmt"
	redsync "github.com/Nghiait123456/redlock"
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
	customExpiry := redsync.OptionFunc(func(mutex *redsync.Mutex) {
		mutex.SetExpiry(2 * time.Second)
	})

	customTries := redsync.OptionFunc(func(mutex *redsync.Mutex) {
		mutex.SetTries(3)
	})

	customDelayFc := redsync.OptionFunc(func(mutex *redsync.Mutex) {
		mutex.SetDelayFunc(func(tries int) time.Duration {
			return 200 * time.Microsecond
		})
	})

	customDiftFactor := redsync.OptionFunc(func(mutex *redsync.Mutex) {
		mutex.SetDriftFactor(0.01)
	})

	mutex := rs.NewMutex("test-redsync2", customExpiry, customTries, customDelayFc, customDiftFactor)
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
