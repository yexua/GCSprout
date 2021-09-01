// @Author : Lik
// @Time   : 2021/5/26
package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	i := a()
	if i == int64(1) {
		fmt.Println(true)
	} else {
		fmt.Println(false)
	}
}

func a() interface{} {
	var i int64 = 1
	return i
}

func TestLock(t *testing.T) {

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			testRedisLock(fmt.Sprint(i))
		}()
	}
	wg.Wait()
}

func testRedisLock(threadName string) {
	var param lockParam
	param.lockKey = "counter_lock"
	param.tryLockTime = time.Second * 2
	param.holdLockTime = time.Second * 10
	opt := &redis.Options{
		Addr:     "120.26.44.207:26379",
		Password: "Ndkj1234",
		DB:       1,
	}
	d := NewDisLockRedis(opt, param)
	defer func() {
		unLock := d.UnLock()
		log.Printf("协程: %s - 释放锁结果：%v\n", threadName, unLock)
	}()
	lockFlag := d.Lock()
	log.Printf("协程: %s - 加锁结果：%v\n", threadName, lockFlag)
	if lockFlag {
		log.Println("正在执行任务......")
		time.Sleep(time.Second * 15)
	}

}
