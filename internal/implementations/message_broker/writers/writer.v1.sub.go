package writers

import (
	"bmt_upload_service/global"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	writer  *kafka.Writer
	once    sync.Once
	closeCh = make(chan struct{})
)

func initKafkaWriter() {
	once.Do(func() {
		writer = &kafka.Writer{
			Addr: kafka.TCP(
				global.Config.ServiceSetting.KafkaSetting.KafkaBroker_1,
				global.Config.ServiceSetting.KafkaSetting.KafkaBroker_2,
				global.Config.ServiceSetting.KafkaSetting.KafkaBroker_3,
			),
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 1 * time.Millisecond,
			MaxAttempts:  3,
			BatchSize:    100,
			WriteTimeout: 5 * time.Second,
		}
	})
}

func ensureTopicExists(topic string) error {
	conn, err := kafka.Dial("tcp", global.Config.ServiceSetting.KafkaSetting.KafkaBroker_1)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err == nil && len(partitions) > 0 {
		// Topic is exists
		return nil
	}

	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
		ConfigEntries:     []kafka.ConfigEntry{},
	})
}
