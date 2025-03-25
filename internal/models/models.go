package models

// PostRerquest - структура запроса на создание короткой ссылки
type PostRerquest struct {
	URL string `json:"url" binding:"required"`
}

// PostResponse - структура ответа на запрос на создание короткой ссылки
type PostResponse struct {
	Result string `json:"result"`
}

// LocalLink - структура хранения ссылок
type LocalLink struct {
	ID  string `json:"short_url"`
	URL string `json:"original_url"`
}

// PostBatchReq - структура запроса на создание нескольких коротких ссылок
type PostBatchReq struct {
	ID  string `json:"correlation_id" binding:"required"`
	URL string `json:"original_url" binding:"required"`
}

// PostBatchResp - структура ответа на запрос на создание нескольких коротких ссылок
type PostBatchResp struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}

// UsersUrlResp - структура ответа на запрос на получение всех коротких ссылок пользователя
type UsersUrlResp struct {
	ID  string `json:"short_url"`
	URL string `json:"original_url"`
}
