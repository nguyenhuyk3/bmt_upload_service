package readers

import (
	"bmt_upload_service/global"
	"bmt_upload_service/internal/implementations/message_broker/writers"
	"bmt_upload_service/internal/services"
	"log"
)

type Reader struct {
	UploadService services.IUpload
	Writer        *writers.Writer
}

func NewReader(
	uploadService services.IUpload,
	writer *writers.Writer,
) *Reader {
	return &Reader{
		UploadService: uploadService,
		Writer:        writer,
	}
}

var topics = []string{
	global.UPLOAD_IMAGE_TOPIC,
	global.UPLOAD_VIDEO_TOPIC,
}

func (r *Reader) InitReaders() {
	log.Println("=============== Upload Service is listening for messages... ===============")

	for _, topic := range topics {
		go r.startReader(topic)
	}
}
