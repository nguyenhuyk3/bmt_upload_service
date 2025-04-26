package upload

import (
	"bmt_upload_service/global"
	"fmt"
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
