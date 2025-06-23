package readers

import (
	"bmt_upload_service/dto/messages"
	"bmt_upload_service/global"
	"bmt_upload_service/internal/services"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSReader struct {
	QueueURL  string
	SQSClient *sqs.Client
	Writer    services.IMessageBrokerWriter
}

const (
	film_image_base_key = "film-images/"
	film_video_base_key = "film-videos/"
	fab_image_base_ket  = "fab-images/"
)

func NewSQSReader(
	queueURL string,
	writer services.IMessageBrokerWriter) *SQSReader {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(global.Config.Server.AWSRegion))
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	return &SQSReader{
		QueueURL:  queueURL,
		SQSClient: sqsClient,
		Writer:    writer,
	}
}

func (sr *SQSReader) IniSQSReader() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		out, err := sr.SQSClient.ReceiveMessage(ctx,
			&sqs.ReceiveMessageInput{
				QueueUrl:            &sr.QueueURL,
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
				VisibilityTimeout:   30,
			})
		cancel()

		if err != nil {
			log.Printf("receive message error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(out.Messages) == 0 {
			continue
		}

		for _, msg := range out.Messages {
			var snsMessage map[string]interface{}
			if err := json.Unmarshal([]byte(*msg.Body), &snsMessage); err != nil {
				log.Printf("invalid JSON in SNS message: %v", err)
				continue
			}

			messageStr, ok := snsMessage["Message"].(string)
			if !ok {
				log.Printf("invalid Message format in SNS notification")
				continue
			}

			var evt messages.S3Event
			if err := json.Unmarshal([]byte(messageStr), &evt); err != nil {
				log.Printf("invalid JSON in S3 event message: %v", err)
				continue
			}

			if len(evt.Records) == 0 {
				log.Printf("empty Records (maybe TestEvent), skipping...")

				// sr.deleteMessage(msg)

				continue
			}

			for _, data := range evt.Records {
				switch data.EventName {
				case global.SQS_PUT_EVENT:
					keyPrefix := extractPrefix(data.S3.Object.Key)

					switch keyPrefix {
					case film_image_base_key, film_video_base_key:
						productId, ext, err := parseKey(data.S3.Object.Key)
						if err != nil {
							log.Printf("parseKey error: %v", err)
							continue
						}

						objectURL := generateObjectURL(
							global.Config.ServiceSetting.S3Setting.FilmBucketName,
							global.Config.ServiceSetting.S3Setting.AwsRegion,
							data.S3.Object.Key,
						)

						var topic string
						switch ext {
						case "jpg":
							topic = global.RETURNED_FILM_IMAGE_OBJECT_KEY_TOPIC
						case "mp4":
							topic = global.RETURNED_FILM_VIDEO_OBJECT_KEY_TOPIC
						default:
							log.Printf("unsupported file extension: %s", ext)
							continue
						}

						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						err = sr.Writer.SendMessage(ctx, topic, productId,
							messages.ReturnedObjectKeyMessage{
								ProductId: productId,
								ObjectKey: objectURL,
							})
						cancel()

						if err != nil {
							log.Printf("error sending message to Kafka (film - %s) with product id (%s): %v", keyPrefix, productId, err)
							continue
						} else {
							log.Printf("send message to kafka sucessfully (film - %s): %s", keyPrefix, productId)

							sr.deleteMessage(msg)
						}
					case fab_image_base_ket:
						productId, _, err := parseKey(data.S3.Object.Key)
						if err != nil {
							log.Printf("parseKey error: %v", err)
							continue
						}
						objectURL := generateObjectURL(
							global.Config.ServiceSetting.S3Setting.FilmBucketName,
							global.Config.ServiceSetting.S3Setting.AwsRegion,
							data.S3.Object.Key,
						)
						topic := global.RETURNED_FAB_IMAGE_OBJECT_KEY_TOPIC

						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						err = sr.Writer.SendMessage(ctx, topic, productId,
							messages.ReturnedObjectKeyMessage{
								ProductId: productId,
								ObjectKey: objectURL,
							})
						cancel()

						if err != nil {
							log.Printf("error sending message to Kafka (fab - %s) with product id (%s): %v", keyPrefix, productId, err)
							continue
						} else {
							log.Printf("send message to kafka sucessfully (fab - %s) with product id: %s", keyPrefix, productId)

							sr.deleteMessage(msg)
						}
					default:
						log.Printf("invalid prefix: %s", keyPrefix)
						continue
					}
				}
			}

			// sr.deleteMessage(msg)
		}
	}
}
