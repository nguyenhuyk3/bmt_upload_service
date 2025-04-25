package messagebroker

import (
	"bmt_upload_service/global"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
			Addr:         kafka.TCP(global.Config.ServiceSetting.KafkaSetting.KafkaBroker_1),
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

func SendMessage(topic string, key string, value interface{}) error {
	if writer == nil {
		initKafkaWriter()
	}

	if err := ensureTopicExists(topic); err != nil {
		log.Printf("failed to ensure topic exists: %v", err)
		return err
	}

	msgBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: msgBytes,
	})
	if err != nil {
		log.Printf("failed to send message to Kafka: %v", err)
		return err
	}

	// log.Printf("message sent to Kafka topic %s", topic)
	return nil
}

func Close() {
	if writer != nil {
		writer.Close()
		log.Println("kafka producer closed")
	}

	close(closeCh)
}
