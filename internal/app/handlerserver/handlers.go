package handlerserver

import (
	"encoding/json"
	"fmt"
	"github.com/Ippolid/shortLink/internal/app"
	"github.com/Ippolid/shortLink/internal/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// PostCreate обрабатывает POST-запрос для создания короткой ссылки.
// Принимает тело запроса с оригинальной ссылкой, возвращает короткую ссылку.
// Требует авторизации, получая userID из контекста.
// Возвращает статус 201 Created при успешном создании, 409 Conflict при дублировании.
// @Summary Создание короткой ссылки
// @Description Принимает оригинальную ссылку и возвращает короткую
// @Tags ссылки
// @Accept plain
// @Produce plain
// @Param url body string true "Оригинальный URL"
// @Success 201 {string} string "Короткая ссылка"
// @Failure 400 {string} string "Ошибка ввода"
// @Failure 401 {string} string "Неавторизованный запрос"
// @Failure 409 {string} string "Ссылка уже существует"
// @Router / [post]
func (s *Server) PostCreate(c *gin.Context) {
	// Получаем user_id из контекста (его устанавливает AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userID.(string)

	// Читаем тело запроса
	val, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Can't read body")
		return
	}

	// Генерируем уникальный short_id
	id := app.GenerateShortID(val)

	// Проверяем наличие базы данных
	if s.Db == nil {
		// Проверяем, есть ли уже такой short_id
		_, exist := s.database.Data[id]
		if exist {
			c.String(http.StatusConflict, s.Adr+id)
			return
		}
		s.database.SaveLink(val, id)
		// Сохраняем ссылку в локальную "базу"
		s.database.SaveUserLink(userIDStr, string(val))
	} else {
		// Сохраняем ссылку в БД (если она есть)
		err = s.Db.InsertLink(id, string(val), userIDStr)
		if err != nil {
			fmt.Println(err)
			if strings.Contains(err.Error(), "link exists") {
				c.String(http.StatusConflict, s.Adr+id)
				return
			}
			c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
			return
		}
	}

	// Устанавливаем content-type и возвращаем результат
	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, s.Adr+id)
}

// GetID обрабатывает GET-запрос для получения оригинальной ссылки по идентификатору.
// Получает ID из URL-параметра и перенаправляет на оригинальную ссылку.
// Возвращает 307 TemporaryRedirect при успехе, 400 BadRequest при ошибке, 410 Gone если ссылка удалена.
func (s *Server) GetID(c *gin.Context) {
	var val string
	var err error
	var exist bool
	id := c.Param("id")
	if s.Db == nil {
		val, exist = s.database.Data[id]
		if !exist {
			c.String(http.StatusBadRequest, "Can't find link")
			return
		}
	} else {
		val, exist, err = s.Db.GetLink(id)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
			return
		}
		if exist {
			c.String(http.StatusGone, "Can't find link")
			return
		}
	}

	fmt.Println(val)
	if err != nil {
		c.String(http.StatusBadRequest, "Can't find link")
		return
	}

	c.Header("content-type", "text/plain")
	c.Redirect(http.StatusTemporaryRedirect, val)
}

// PingDB проверяет доступность базы данных.
// Возвращает 200 OK если база доступна, 500 InternalServerError в противном случае.
func (s *Server) PingDB(c *gin.Context) {
	b, err := s.Db.Ping()
	if err != nil {
		c.String(http.StatusInternalServerError, "DB is not available")
		return
	}
	if b {
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusInternalServerError)
}

// PostAPI обрабатывает POST-запрос к API для создания короткой ссылки.
// Принимает JSON-запрос с полем URL, возвращает JSON-ответ с короткой ссылкой.
// Требует авторизации, получая userID из контекста.
// Возвращает статус 201 Created при успешном создании, 409 Conflict при дублировании.
func (s *Server) PostAPI(c *gin.Context) {
	var req models.PostRerquest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userID.(string)

	id := app.GenerateShortID([]byte(req.URL))
	if s.Db == nil {
		_, exist := s.database.Data[id]
		if exist {
			response := models.PostResponse{
				Result: s.Adr + id,
			}
			c.JSON(http.StatusConflict, response)
			return
		}
		s.database.SaveLink([]byte(req.URL), id)
		s.database.SaveUserLink(userIDStr, req.URL)
	} else {
		err := s.Db.InsertLink(id, req.URL, userIDStr)
		if err != nil {
			fmt.Println(err)
			if strings.Contains(err.Error(), "link exists") {
				response := models.PostResponse{
					Result: s.Adr + id,
				}
				c.JSON(http.StatusConflict, response)
				return
			}
			c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
			return
		}
	}
	response := models.PostResponse{
		Result: s.Adr + id,
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
}

// PostBatch обрабатывает пакетный POST-запрос для создания нескольких коротких ссылок.
// Принимает JSON-массив с парами {ID, URL}, возвращает JSON-массив с короткими ссылками.
// Требует авторизации, получая userID из контекста.
// Возвращает статус 201 Created при успешном создании.
func (s *Server) PostBatch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userID.(string)

	var req []models.PostBatchReq
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}
	var otv models.PostBatchResp
	var resp []models.PostBatchResp
	for _, r := range req {
		if r.ID != "" && r.URL != "" {
			otv.ID = r.ID
			k := app.GenerateShortID([]byte(r.URL))
			otv.URL = s.Adr + k

			if s.Db == nil {
				s.database.SaveLink([]byte(r.URL), k)
				s.database.SaveUserLink(userIDStr, r.URL)
			} else {
				err := s.Db.InsertLink(k, r.URL, userIDStr)
				if err != nil {
					c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
				}

			}
			resp = append(resp, otv)
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, resp)

}

// GetUserURLs возвращает список всех ссылок пользователя.
// Требует авторизации, получая userID из контекста.
// Возвращает JSON-массив с короткими и оригинальными ссылками.
// Возвращает статус 200 OK при успешном получении, 204 NoContent если ссылок нет.
//
// @Summary Получение всех ссылок пользователя
// @Description Возвращает список всех ссылок, созданных пользователем
// @Tags пользователь
// @Produce json
// @Success 200 {array} models.UsersUrlResp "Список ссылок пользователя"
// @Success 204 {string} string "Ссылки не найдены"
// @Failure 401 {string} string "Неавторизованный запрос"
// @Router /api/user/urls [get]
func (s *Server) GetUserURLs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var userURLs = make([]string, 0)
	var found bool
	var err error

	userIDStr := userID.(string)
	if s.Db == nil {
		userURLs, found = s.database.LoadUserLink(userIDStr)

		if !found || len(userURLs) == 0 {
			c.Status(http.StatusNoContent)
			return
		}
	} else {
		userURLs, err = s.Db.GetLinksByUserID(userIDStr)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
		}
		fmt.Println(userURLs)
		if len(userURLs) == 0 {
			c.Status(http.StatusNoContent)
			return
		}

	}
	fmt.Println(userURLs)
	var otv models.UsersUrlResp
	var resp []models.UsersUrlResp
	var shortlink string

	for _, r := range userURLs {
		id := app.GenerateShortID([]byte(r))
		shortlink = s.Adr + id
		otv.ID = shortlink
		otv.URL = r
		resp = append(resp, otv)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resp)
}

// DeleteLinks помечает указанные ссылки как удаленные.
// Принимает JSON-массив с идентификаторами ссылок.
// Удаляет только ссылки, принадлежащие текущему пользователю.
// Требует авторизации, получая userID из контекста.
// Возвращает статус 202 Accepted при успешном выполнении.
func (s *Server) DeleteLinks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var req []string
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}
	s.Db.Dellink(req, userID.(string))
	c.Status(http.StatusAccepted)
}
