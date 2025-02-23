package models

type PostRerquest struct {
	Url string `json:"url" binding:"required"`
}
type PostResponse struct {
	Result string `json:"result"`
}
