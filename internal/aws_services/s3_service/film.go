package s3service

import (
	"bmt_upload_service/dto/messages"
	"bmt_upload_service/global"
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	film_image_base_key = "film-images/"
	film_video_base_key = "film-videos/"
)

func UploadFilmImageToS3(message messages.UploadImageMessage) error {
	file, err := os.Open(message.ImageUrl)
	if err != nil {
		return fmt.Errorf("cann't open file: %v", err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("reading file err: %v", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(global.Config.ServiceSetting.S3Setting.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			global.Config.ServiceSetting.S3Setting.AwsAccessKeyId,
			global.Config.ServiceSetting.S3Setting.AwsSercetAccessKeyId,
			""),
	})
	if err != nil {
		return fmt.Errorf("connecting AWS err: %v", err)
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:         aws.String(film_image_base_key + message.ImageUrl),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String("image/jpeg"),
	})

	if err != nil {
		return fmt.Errorf("upload to S3 (image) failure: %v", err)
	}

	return nil
}

func UploadFilmVideoToS3(message messages.UploadVideoMessage) error {
	file, err := os.Open(message.VideoUrl)
	if err != nil {
		return fmt.Errorf("cann't open file: %v", err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("reading file err: %v", err)
	}

	// Lấy phần mở rộng file (ví dụ: ".mp4")
	ext := filepath.Ext(message.VideoUrl)
	// Xác định MIME type (ví dụ: "video/mp4")
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(global.Config.ServiceSetting.S3Setting.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			global.Config.ServiceSetting.S3Setting.AwsAccessKeyId,
			global.Config.ServiceSetting.S3Setting.AwsSercetAccessKeyId,
			""),
	})
	if err != nil {
		return fmt.Errorf("connecting AWS err: %v", err)
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:         aws.String(film_video_base_key + message.VideoUrl),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(mimeType),
	})

	if err != nil {
		return fmt.Errorf("upload to S3 (video) failure: %v", err)
	}

	return nil
}
