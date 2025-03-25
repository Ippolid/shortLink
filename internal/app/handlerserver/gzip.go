package handlerserver

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// gzipMiddleware - Middleware для сжатия ответов
func gzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем поддержку gzip у клиента
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		contentType := c.Writer.Header().Get("Content-Type")
		if !(strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")) {
			return
		}

		// Выполняем запрос и перехватываем ответ
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Set("Vary", "Accept-Encoding")

		// Создаём gzip.Writer поверх ResponseWriter
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		// Оборачиваем ResponseWriter для сжатия данных
		gzWriter := &gzipResponseWriter{ResponseWriter: c.Writer, Writer: gz}
		c.Writer = gzWriter

		c.Next()
	}
}

// gzipResponseWriter - Обёртка для сжатия ответа
type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

// Write - Реализация метода Write для gzipResponseWriter
func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

// gzipDecompressMiddleware - Middleware для распаковки входящего запроса
func gzipDecompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, пришёл ли сжатый запрос
		if c.GetHeader("Content-Encoding") != "gzip" {
			c.Next()
			return
		}

		// Распаковываем тело запроса
		gz, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to decompress request"})
			return
		}
		defer gz.Close()

		// Читаем распакованные данные в буфер
		var body bytes.Buffer
		_, err = io.Copy(&body, gz)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to read decompressed request"})
			return
		}

		// Заменяем тело запроса на распакованные данные
		c.Request.Body = io.NopCloser(&body)
		c.Next()
	}
}
