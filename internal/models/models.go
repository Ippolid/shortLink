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

type PostBatchReq struct {
	ID  string `json:"correlation_id" binding:"required"`
	URL string `json:"original_url" binding:"required"`
}

type PostBatchResp struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}

type GETUserLinks struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type UserCookie struct {
	UserID string `json:"user_id"`
	Sign   string `json:"sign"`
}
