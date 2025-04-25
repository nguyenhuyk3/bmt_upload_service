package messagebroker

import (
	"bmt_upload_service/dto/messages"
	"bmt_upload_service/global"
	s3service "bmt_upload_service/internal/aws_services/s3_service"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var topics = []string{
	global.UPLOAD_IMAGE_TOPIC,
	global.UPLOAD_VIDEO_TOPIC,
}

func InitReaders() {
	log.Println("=============== Upload Service is listening for messages... ===============")

	for _, topic := range topics {
		go startReader(topic)
	}
}

func startReader(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{global.Config.ServiceSetting.KafkaSetting.KafkaBroker_1},
		GroupID:        global.UPLOAD_SERVICE_GROUP,
		Topic:          topic,
		CommitInterval: time.Second * 5,
	})
	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		processMessage(topic, message.Value)
	}
}

func processMessage(topic string, value []byte) {
	switch topic {
	case global.UPLOAD_IMAGE_TOPIC:
		var uploadMessage messages.UploadImageMessage
		if err := json.Unmarshal(value, &uploadMessage); err != nil {
			log.Printf("failed to unmarshal image message: %v\n", err)
			return
		}

		handleImageUpload(uploadMessage)

	case global.UPLOAD_VIDEO_TOPIC:
		var uploadMessage messages.UploadVideoMessage
		if err := json.Unmarshal(value, &uploadMessage); err != nil {
			log.Printf("failed to unmarshal video message: %v\n", err)
			return
		}

		handleVideoUpload(uploadMessage)

	default:
		log.Printf("unknown topic received: %s\n", topic)
	}
}

func handleImageUpload(message messages.UploadImageMessage) {
	objectKey, err := s3service.UploadFilmImageToS3(message)
	if err != nil {
		log.Printf("failed to upload image: %v\n", err)
	} else {
		log.Printf("successfully uploaded image for ProductID: %s\n", message.ProductId)
		// Retry send message to topic with 3 times
		topic := global.RETURNED_IMAGE_OBJECT_KEY_TOPIC
		err = sendReturnedObjectKey(topic, message.ProductId, objectKey)
		if err != nil {
			log.Printf("failed to send message to Kafka (%s): %\n", topic, err)
		} else {
			log.Printf("send message to (%s) topic\n", topic)
		}
	}
}

func handleVideoUpload(message messages.UploadVideoMessage) {
	objectKey, err := s3service.UploadFilmVideoToS3(message)
	if err != nil {
		log.Printf("failed to upload video: %v\n", err)
	} else {
		log.Printf("successfully uploaded video for ProductID: %s\n", message.ProductId)
		// Retry send message to topic with 3 times
		topic := global.RETURNED_VIDEO_OBJECT_KEY_TOPIC
		err = sendReturnedObjectKey(topic, message.ProductId, objectKey)
		if err != nil {
			log.Printf("failed to send message to Kafka (%s): %v\n", topic, err)
		} else {
			log.Printf("send message to (%s) topic\n", topic)
		}
	}
}

func sendReturnedObjectKey(topic,
	productId, objectKey string) error {
	return SendMessage(topic,
		topic,
		messages.ReturnedObjectKeyMessage{
			ProductId: productId,
			ObjectKey: objectKey})
}
