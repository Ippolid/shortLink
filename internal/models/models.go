package models

type PostRerquest struct {
	URL string `json:"url" binding:"required"`
}
type PostResponse struct {
	Result string `json:"result"`
}

type LocalLink struct {
	ID  string `json:"short_url"`
	URL string `json:"original_url"`
}
