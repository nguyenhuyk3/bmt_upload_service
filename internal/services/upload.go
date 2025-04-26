package services

import "bmt_upload_service/dto/messages"

type IUpload interface {
	UploadFilmImageToS3(message messages.UploadImageMessage) (string, error)
	UploadFilmVideoToS3(message messages.UploadVideoMessage) (string, error)
}
