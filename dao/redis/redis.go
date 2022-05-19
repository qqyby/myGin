package redis

import (
	"context"
	"time"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func initClient(masterName string, addr []string, pwd string, db, poolSize int) {
	if len(addr) == 0 {
		panic("addr empty")
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr[0],
		Password: pwd,
		DB:       db,
		PoolSize: poolSize,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

func Init(masterName string, addr []string, pwd string, db, poolSize int) error {
	// 用于内部测试时候使用redis的单节点启动
	if settings.AppCfg.RunMode == global.InnerTestModel {
		initClient(masterName, addr, pwd, db, poolSize)
		return nil
	}

	rdb = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addr,
		Password:      pwd, // 没有设置密码
		DB:            db,  // 使用默认第一个db
		PoolSize:      poolSize,
	})

	_, err := rdb.Ping(context.Background()).Result()
	return err
}

func Close() error {
	return rdb.Close()
}

// redis get string 的通用操作
func Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), cmdTimeOut)
	defer cancel()
	rest, err := rdb.Get(ctx, key).Result()
	return rest, err
}

// redis set string 的通用操作
func Set(key, value string, t time.Duration) error {
	ctx, cancel := context.WithTimeout(context.TODO(), cmdTimeOut)
	defer cancel()
	return rdb.Set(ctx, key, value, t).Err()
}

// redis set string 的通用操作
func Del(key string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), cmdTimeOut)
	defer cancel()
	return rdb.Del(ctx, key).Err()
}

// 分布式锁的简单实现
func Lock(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	b, err := rdb.SetNX(ctx, key, 1, 10*time.Second).Result()
	return b, err
}

// 删除锁
func UnLock(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	nums, err := rdb.Del(ctx, key).Result()
	return nums, err
}
