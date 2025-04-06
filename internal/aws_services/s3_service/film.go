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
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	film_image_base_key = "film-images/"
	film_video_base_key = "film-videos/"
)

var (
	once sync.Once
	// s3Session *session.Session
	s3Client *s3.S3
)

func initS3Client() {
	once.Do(func() {
		s3Session, err := session.NewSession(&aws.Config{
			Region: aws.String(global.Config.ServiceSetting.S3Setting.AwsRegion),
			Credentials: credentials.NewStaticCredentials(
				global.Config.ServiceSetting.S3Setting.AwsAccessKeyId,
				global.Config.ServiceSetting.S3Setting.AwsSercetAccessKeyId,
				""),
		})
		if err != nil {
			panic(fmt.Sprintf("failed to create AWS session: %v", err))
		}

		s3Client = s3.New(s3Session)
	})
}

func UploadFilmImageToS3(message messages.UploadImageMessage) error {
	if s3Client == nil {
		initS3Client()
	}

	file, err := os.Open(message.ImageUrl)
	if err != nil {
		return fmt.Errorf("can't open file: %v", err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("reading file err: %v", err)
	}

	objectKey := film_image_base_key + filepath.Base(message.ImageUrl)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return fmt.Errorf("upload to S3 (image) failure: %v", err)
	}

	return nil
}

func UploadFilmVideoToS3(message messages.UploadVideoMessage) error {
	if s3Client == nil {
		initS3Client()
	}

	file, err := os.Open(message.VideoUrl)
	if err != nil {
		return fmt.Errorf("can't open file: %v", err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("reading file err: %v", err)
	}
	// Get the file extension (eg ".mp4")
	ext := filepath.Ext(message.VideoUrl)
	// Specify the MIME type (eg "video/mp4")
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	objectKey := film_video_base_key + filepath.Base(message.VideoUrl)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(global.Config.ServiceSetting.S3Setting.FilmBucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return fmt.Errorf("upload to S3 (video) failure: %v", err)
	}

	return nil
}
