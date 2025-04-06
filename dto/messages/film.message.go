package messages

// Incomming
type UploadImageMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	ImageUrl  string `json:"image_url" binding:"required"`
}

type UploadVideoMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	VideoUrl  string `json:"video_url" binding:"required"`
}

// Sending
type ReturnedObjectKeyMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	ObjectKey string `json:"object_key" binding:"required"`
}
