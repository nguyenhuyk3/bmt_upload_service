package messages

type ReturnedObjectKeyMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	ObjectKey string `json:"object_key" binding:"required"`
}

// From SQS
type S3EventRecordMessage struct {
	EventName string `json:"eventName"`
	S3        struct {
		Bucket struct {
			Name string
		} `json:"bucket"`
		Object struct {
			Key string `json:"key"`
		} `json:"object"`
	} `json:"s3"`
}

type S3Event struct {
	Records []S3EventRecordMessage `json:"Records"`
}
