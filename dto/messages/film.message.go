package messages

type UploadImageMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	ImageUrl  string `json:"image_url" binding:"required"`
}

type UploadVideoMessage struct {
	ProductId string `json:"product_id" binding:"required"`
	VideoUrl  string `json:"video_url" binding:"required"`
}
