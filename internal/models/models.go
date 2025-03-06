package models

import "github.com/golang-jwt/jwt/v4"

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

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type GETUserLinks struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
