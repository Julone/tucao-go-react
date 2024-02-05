package middleware

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/segmentio/kafka-go"
	"time"
	"tuxiaocao/pkg/logger"
)

func HookupFromKafka() {
	// 指定要连接的topic和partition
	topic := "my-topic"
	partition := 0
	fmt.Println("Kafka is listening")
	// 连接至Kafka的leader节点
	conn := kafka.NewReader(kafka.ReaderConfig{
		Brokers:                []string{"localhost:29092"},
		Topic:                  topic,
		MinBytes:               10e3,
		MaxBytes:               10e6,
		Partition:              partition,
		MaxAttempts:            3,
		HeartbeatInterval:      time.Second,
		PartitionWatchInterval: time.Second,
		Logger: kafka.LoggerFunc(func(msg string, a ...interface{}) {
			logger.Log.Infof(msg, a...)
		}),
	})
	err := conn.SetOffset(80)
	if err != nil {
		logger.Log.Info(err)
		return
	}

	for {
		v, err := conn.FetchMessage(context.Background())
		if err != nil {
			logger.Log.Info(err)
			break
		}
		log.Info("kafka message: ", v.Offset, string(v.Value), v.Headers)
	}

	defer func() {
		conn.Close().Error()
	}()
}
