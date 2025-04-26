package upload

import (
	"bmt_upload_service/dto/messages"
	"bmt_upload_service/global"
	"bmt_upload_service/internal/services"
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type UploadService struct {
}

// UploadFilmImageToS3 implements services.IS3.
func (us *UploadService) UploadFilmImageToS3(message messages.UploadImageMessage) (string, error) {
	if s3Client == nil {
		initS3Client()
	}

	file, err := os.Open(message.ImageUrl)
	if err != nil {
		return "", fmt.Errorf("can't open file: %v", err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("reading file err: %v", err)
	}

	ext := filepath.Ext(message.ImageUrl)
	if ext == "" {
		ext = ".jpg"
	}

	newFileName := uuid.New().String() + ext

	objectKey := film_image_base_key + newFileName
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", fmt.Errorf("upload to S3 (image) failure: %v", err)
	}

	return fmt.Sprintf("https://%s.%s/%s",
		global.Config.ServiceSetting.S3Setting.FilmBucketName,
		global.Config.ServiceSetting.S3Setting.AwsRegion,
		objectKey,
	), nil
}

// UploadFilmVideoToS3 implements services.IS3.
func (us *UploadService) UploadFilmVideoToS3(message messages.UploadVideoMessage) (string, error) {
	if s3Client == nil {
		initS3Client()
	}

	file, err := os.Open(message.VideoUrl)
	if err != nil {
		return "", fmt.Errorf("can't open file: %v", err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("reading file err: %v", err)
	}
	// Get the file extension (eg ".mp4")
	ext := filepath.Ext(message.VideoUrl)
	// Specify the MIME type (eg "video/mp4")
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	newFileName := uuid.New().String() + ext
	objectKey := film_video_base_key + newFileName

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", fmt.Errorf("upload to S3 (video) failure: %v", err)
	}

	return fmt.Sprintf("https://%s.%s/%s",
		global.Config.ServiceSetting.S3Setting.FilmBucketName,
		global.Config.ServiceSetting.S3Setting.AwsRegion,
		objectKey,
	), nil
}

func NewUploadService() services.IUpload {
	return &UploadService{}
}
