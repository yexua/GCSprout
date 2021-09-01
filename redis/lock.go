// Package redis
// @Author : Lik
// @Time   : 2021/5/25
package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"strconv"
	"time"
)

//var ctx = context.Background()
//var lockKey = "counter_lock"
//var redisLockValue string
var unLockScript = "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
var unLockSuccess int64 = 1

type DisLockRedis struct {
	redisLockKey string
	redisLockVal string
	// 尝试获取锁的截至时间
	tryLockEndTime time.Time
	ctx            context.Context
	cancelFunc     context.CancelFunc
	redis          *redis.Client
	lockParam
}

type lockParam struct {
	//锁的key
	lockKey string
	//尝试获得锁的时间
	tryLockTime time.Duration
	//尝试获得锁后，持有锁的时间
	holdLockTime time.Duration
}

func NewDisLockRedis(opt *redis.Options, param lockParam) *DisLockRedis {
	conn := NewRedisConn(opt)
	//uuid, err := exec.Command("uuidgen").Output()
	//if err != nil {
	//
	//}
	rand.Seed(time.Now().Unix())

	this := &DisLockRedis{
		redisLockKey:   param.lockKey,
		redisLockVal:   strconv.Itoa(rand.Int()),
		tryLockEndTime: time.Now().Add(param.tryLockTime),
		redis:          conn,
		ctx:            context.Background(),
		lockParam:      param,
	}
	return this
}

func NewRedisConn(opt *redis.Options) *redis.Client {

	rdb := redis.NewClient(opt)
	return rdb
}

func (d *DisLockRedis) Lock() bool {
	holdLockSuccess := false
	for {
		holdLockSuccess = d.TryLock()
		if holdLockSuccess {
			break
		}
		time.Sleep(50)
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	d.cancelFunc = cancelFunc

	if holdLockSuccess {
		d.renew(ctx)
	}

	return holdLockSuccess
}
func (d *DisLockRedis) TryLock() bool {
	lockSuccess, err := d.redis.SetNX(d.ctx, d.redisLockKey, d.redisLockVal, d.lockParam.holdLockTime).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return lockSuccess
}
func (d *DisLockRedis) UnLock() bool {
	result, err := d.redis.Eval(d.ctx, unLockScript, []string{d.redisLockKey}, d.redisLockVal).Result()
	//unlock, err := rdb.Del(ctx, key).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	if unLockSuccess == result {
		// 取消续租
		d.cancelFunc()
		return true
	}
	return false
}
func (d *DisLockRedis) renew(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := d.redis.Do(ctx, "EXPIRE", d.redisLockKey, 10).Result()
				if err != nil {
					//log.Printf("reset ttl failure redisLockKey:%s\n", d.redisLockKey)
				}
				if result == unLockSuccess {
					log.Printf("reset ttl success redisLockKey:%s\n", d.redisLockKey)
				}
				if result != unLockSuccess {
					//log.Printf("reset ttl failure redisLockKey:%s\n, result:%v", d.redisLockKey, result)
				}
			}
			time.Sleep(d.holdLockTime / 3)
		}
	}()
}
