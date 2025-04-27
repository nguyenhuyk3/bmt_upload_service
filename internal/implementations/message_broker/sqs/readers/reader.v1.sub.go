package readers

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func parseKey(objectKey string) (productId string, ext string, err error) {
	// 1. Get filename (remove folder)
	filename := path.Base(objectKey) // "11-3f0b078b-9bac-44d3-8c96-d2cadba394bf.mp4"
	// 2. Seperate file tail
	dotIdx := strings.LastIndex(filename, ".")
	if dotIdx == -1 {
		return "", "", fmt.Errorf("invalid objectKey: no file extension (%s)", objectKey)
	}
	ext = filename[dotIdx+1:]           // "mp4"
	nameWithoutExt := filename[:dotIdx] // "11-3f0b078b-9bac-44d3-8c96-d2cadba394bf"
	// 3. Get productId
	dashIdx := strings.Index(nameWithoutExt, "-")
	if dashIdx == -1 {
		return "", "", fmt.Errorf("invalid objectKey: no dash in filename (%s)", objectKey)
	}
	productId = nameWithoutExt[:dashIdx] // "11"

	return productId, ext, nil
}

func generateObjectURL(bucketName, region, objectKey string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey)
}

func (sr *SQSReader) deleteMessage(msg sqstypes.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := sr.SQSClient.DeleteMessage(ctx,
		&sqs.DeleteMessageInput{
			QueueUrl:      &sr.QueueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
	if err != nil {
		log.Printf("delete message error (SQS): %v", err)
	}
}
