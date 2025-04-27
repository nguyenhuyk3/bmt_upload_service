package initializations

import (
	"bmt_upload_service/global"
	"bmt_upload_service/internal/implementations/message_broker/kafka/writers"
	"bmt_upload_service/internal/implementations/message_broker/sqs/readers"
)

func initSQSQueue() {
	kafkaWriter := writers.NewKafkaWriter()
	sqsReader := readers.NewSQSReader(global.Config.Server.SQSUrl, kafkaWriter)

	sqsReader.IniSQSReader()

	select {}
}
