package readers

import (
	"bmt_upload_service/dto/messages"
	"bmt_upload_service/global"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func (r *Reader) startReader(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_1,
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_2,
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_3,
		},
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

		r.processMessage(topic, message.Value)
	}
}

func (r *Reader) processMessage(topic string, value []byte) {
	switch topic {
	// Handle uploading image
	case global.UPLOAD_IMAGE_TOPIC:
		var uploadMessage messages.UploadImageMessage
		if err := json.Unmarshal(value, &uploadMessage); err != nil {
			log.Printf("failed to unmarshal image message: %v\n", err)
			return
		}

		r.handleImageUpload(uploadMessage)

	// Handle uploading video
	case global.UPLOAD_VIDEO_TOPIC:
		var uploadMessage messages.UploadVideoMessage
		if err := json.Unmarshal(value, &uploadMessage); err != nil {
			log.Printf("failed to unmarshal video message: %v\n", err)
			return
		}

		r.handleVideoUpload(uploadMessage)

	default:
		log.Printf("unknown topic received: %s\n", topic)
	}
}

func (r *Reader) handleImageUpload(message messages.UploadImageMessage) {
	objectKey, err := r.UploadService.UploadFilmImageToS3(message)
	if err != nil {
		log.Printf("failed to upload image: %v\n", err)
	} else {
		log.Printf("successfully uploaded image for ProductID: %s\n", message.ProductId)

		topic := global.RETURNED_IMAGE_OBJECT_KEY_TOPIC
		err = r.sendReturnedObjectKey(topic, message.ProductId, objectKey)
		if err != nil {
			log.Printf("failed to send message to Kafka (%s): %\n", topic, err)
		} else {
			log.Printf("send message to (%s) topic\n", topic)
		}
	}
}

func (r *Reader) handleVideoUpload(message messages.UploadVideoMessage) {
	objectKey, err := r.UploadService.UploadFilmVideoToS3(message)
	if err != nil {
		log.Printf("failed to upload video: %v\n", err)
	} else {
		log.Printf("successfully uploaded video for ProductID: %s\n", message.ProductId)

		topic := global.RETURNED_VIDEO_OBJECT_KEY_TOPIC
		err = r.sendReturnedObjectKey(topic, message.ProductId, objectKey)
		if err != nil {
			log.Printf("failed to send message to Kafka (%s): %v\n", topic, err)
		} else {
			log.Printf("send message to (%s) topic\n", topic)
		}
	}
}

func (r *Reader) sendReturnedObjectKey(topic,
	productId, objectKey string) error {
	return r.Writer.SendMessage(topic,
		topic,
		messages.ReturnedObjectKeyMessage{
			ProductId: productId,
			ObjectKey: objectKey})
}
