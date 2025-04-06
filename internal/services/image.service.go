package services

import (
	"bmt_upload_service/dto/request"
)

type IImage interface {
	UploadFilmImageToS3(arg request.UploadFilmImageReq) error
}
