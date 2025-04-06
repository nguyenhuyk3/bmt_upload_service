package request

type UploadFilmImageReq struct {
	ProductId string `json:"product_id" binding:"required"`
	ImageUrl  string `json:"image_url" binding:"required"`
}
