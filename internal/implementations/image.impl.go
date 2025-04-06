package implementations

import (
	"bmt_upload_service/dto/request"
	"bmt_upload_service/internal/services"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type imageService struct {
	S3Client *s3.Client
}

func NewImageService(s3Client *s3.Client) services.IImage {
	return &imageService{S3Client: s3Client}
}

// UploadFilmImageToS3 implements services.IImage.
func (i *imageService) UploadFilmImageToS3(arg request.UploadFilmImageReq) error {
	// // Open file from image url
	// file, err := os.Open(arg.ImageUrl)
	// if err != nil {
	// 	return fmt.Errorf("%v", err)
	// }
	// defer file.Close()

	// s3Key := "products/" + arg.ImageUrl + "/image.jpg"
	// _, err = i.S3Client.PutObject(&s3.PutObjectInput{
	// 	Bucket: aws.String("your-bucket-name"), // Thay bằng bucket của bạn
	// 	Key:    aws.String(s3Key),
	// 	Body:   file,
	// })
	// if err != nil {
	// 	return err
	// }

	// os.Remove(arg.ImageUrl)

	return nil
}
