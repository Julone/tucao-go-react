package setup

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"strconv"
	"time"
	"tuxiaocao/pkg/logger"
	"tuxiaocao/pkg/platform/database"
	redis2 "tuxiaocao/pkg/platform/redis"
	models2 "tuxiaocao/routes/models"
)

func ReadyDatabase() {
	// database init
	_, err := database.OpenDBConnection()
	if err != nil {
		logger.Log.Panicf("mysql is error %v", err)
	}

	err = database.DB.AutoMigrate(models2.Product{}, models2.User{}, models2.LogRecord{})
	if err != nil {
		logger.Log.Errorf("mysql migrate is error %v", err)
	}

}

func ReadyRedisConnection() {
	dbi, _ := strconv.Atoi(os.Getenv("redis_db_number"))
	da := flag.Int("db", dbi, "db init number")

	rds := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       *da,
	})
	redis2.Rds = rds
	result, err := redis2.Rds.Ping(context.Background()).Result()
	if err != nil {
		logger.Log.Errorf("redis pong!", err)
		fmt.Errorf("redis pong!", err)
	}
	logger.Log.Info("redis result ", result)
	rds.Set(context.Background(), "asdfa", time.Now().Format(time.DateTime), 0)
}

func ReadyLogger() {
	logger.Replace(logger.New("./", "rmp", "debug", time.Hour*240, time.Hour*24, 100))
	if logger.Log == nil {
		panic(fmt.Errorf("log start error"))
	}
}

func ReadyEtcd() {
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	_, err := cli.Put(context.Background(), "julone", "lee")
	if err != nil {
		fmt.Printf(err.Error())
	}
	result, _ := cli.Get(context.Background(), "julone")
	fmt.Printf("%v", result.Kvs[0].Value)
}
func InitAll(Components ...string) {
	if len(Components) == 0 {
		Components = []string{"logger", "mysql", "redis", "etcd"}
	}
	defer func() {
		if d := recover(); d != nil {
			fmt.Printf("recover error", d)
		}
	}()
	for _, val := range Components {
		switch val {
		case "logger":
			ReadyLogger()
		case "mysql":
			ReadyDatabase()
		case "redis":
			ReadyRedisConnection()
		case "etcd":
			ReadyEtcd()
		default:
			fmt.Printf("all ready")
		}
	}
}
