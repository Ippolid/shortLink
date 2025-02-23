package models

type PostRerquest struct {
	URL string `json:"url" binding:"required"`
}
type PostResponse struct {
	Result string `json:"result"`
}
