package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddImage(c *gin.Context) {
	planetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}
	// Чтение изображения из запроса
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка чтения фотографии"})
		return
	}
	// Чтение содержимого изображения в байтах
	file, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения фотографии"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения фотографии"})
		return
	}
	// Получение Content-Type из заголовков запроса
	contentType := image.Header.Get("Content-Type")
	// Вызов функции репозитория для добавления изображения
	err = h.Repo.AddPlanetImage(planetId, imageBytes, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Картинка добавлена успешно"})
}
