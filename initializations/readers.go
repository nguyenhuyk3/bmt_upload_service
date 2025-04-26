package initializations

import (
	"bmt_upload_service/internal/implementations/message_broker/readers"
	"bmt_upload_service/internal/implementations/message_broker/writers"
	"bmt_upload_service/internal/implementations/upload"
)

func initReaders() {
	uploadService := upload.NewUploadService()
	writer := writers.NewWriter()
	reader := readers.NewReader(uploadService, writer)

	reader.InitReaders()
}
