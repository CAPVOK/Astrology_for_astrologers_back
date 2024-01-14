package delivery

import (
	"io"
	"net/http"
	"strconv"

	"space/internal/model"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка планеты
// @Description Возращает список всех активных багажей
// @Tags Планета
// @Produce json
// @Param searchCode query string false "Код планеты" Format(email)
// @Success 200 {object} model.PlanetsGetResponse "Список багажей"
// @Failure 500 {object} model.PlanetsGetResponse "Ошибка сервера"
// @Router /planet [get]
func (h *Handler) GetPlanets(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте пп"})
		return
	}
	userID := ctxUserID.(uint)
	searchCode := c.DefaultQuery("searchCode", "")

	planets, err := h.UseCase.GetPlanets(searchCode, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planets": planets.Planets, "constellationID": planets.ConstellationID})
}

// @Summary Получение планеты по ID
// @Description Возвращает информацию о багаже по его ID
// @Tags Планета
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Success 200 {object} model.Planet "Информация о багаже"
// @Failure 400 {object} model.Planet "Некорректный запрос"
// @Failure 500 {object} model.Planet "Внутренняя ошибка сервера"
// @Router /planet/{planet_id} [get]
func (h *Handler) GetPlanetByID(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
		return
	}

	planet, err := h.UseCase.GetPlanetByID(uint(planetID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": planet})
}

// @Summary Создание нового планеты
// @Description Создает новый планета с предоставленными данными
// @Tags Планета
// @Accept json
// @Produce json
// @Param searchCode query string false "Код планеты" Format(email)
// @Param planet body model.PlanetRequest true "Пользовательский объект в формате JSON"
// @Success 200 {object} model.PlanetsGetResponse "Список багажей"
// @Failure 400 {object} model.PlanetsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.PlanetsGetResponse "Внутренняя ошибка сервера"
// @Router /planet/create [post]
func (h *Handler) CreatePlanet(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	searchCode := c.DefaultQuery("searchCode", "")

	var planet model.PlanetRequest

	if err := c.BindJSON(&planet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
		return
	}

	err := h.UseCase.CreatePlanet(userID, planet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	planets, err := h.UseCase.GetPlanets(searchCode, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planets": planets.Planets, "constellationID": planets.ConstellationID})
}

// @Summary Удаление планеты
// @Description Удаляет планета по его ID
// @Tags Планета
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Param searchCode query string false "Код планеты" Format(email)
// @Success 200 {object} model.PlanetsGetResponse "Список багажей"
// @Failure 400 {object} model.PlanetsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.PlanetsGetResponse "Внутренняя ошибка сервера"
// @Router /planet/{planet_id}/delete [delete]
func (h *Handler) DeletePlanet(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	searchCode := c.DefaultQuery("searchName", "")

	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
		return
	}
	err = h.UseCase.DeletePlanet(uint(planetID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	planets, err := h.UseCase.GetPlanets(searchCode, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"planets": planets.Planets, "constellationID": planets.ConstellationID})
}

// @Summary Обновление информации о багаже
// @Description Обновляет информацию о багаже по его ID
// @Tags Планета
// @Accept json
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Success 200 {object} model.Planet "Информация о багаже"
// @Failure 400 {object} model.Planet "Некорректный запрос"
// @Failure 500 {object} model.Planet "Внутренняя ошибка сервера"
// @Router /planet/{planet_id}/update [put]
func (h *Handler) UpdatePlanet(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "недопустимый ИД планеты"}})
		return
	}

	var planet model.PlanetRequest
	if err := c.BindJSON(&planet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
		return
	}

	err = h.UseCase.UpdatePlanet(uint(planetID), uint(userID), planet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedPlanet, err := h.UseCase.GetPlanetByID(uint(planetID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": updatedPlanet})
}

// @Summary Добавление планеты к доставке
// @Description Добавляет планета к доставке по его ID
// @Tags Планета
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Param searchCode query string false "Код планеты" Format(email)
// @Success 200 {object} model.PlanetsGetResponse  "Список багажей"
// @Failure 400 {object} model.PlanetsGetResponse  "Некорректный запрос"
// @Failure 500 {object} model.PlanetsGetResponse  "Внутренняя ошибка сервера"
// @Router /planet/{planet_id}/constellation [post]
func (h *Handler) AddPlanetToConstellation(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	searchCode := c.DefaultQuery("searchCode", "")

	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
		return
	}
	err = h.UseCase.AddPlanetToConstellation(uint(planetID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	planets, err := h.UseCase.GetPlanets(searchCode, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"planets": planets.Planets, "constellationID": planets.ConstellationID})
}

// @Summary Удаление планеты из созвездия
// @Description Удаляет планета из созвездия по его ID
// @Tags Планета
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Param searchCode query string false "Код планеты" Format(email)
// @Success 200 {object} model.PlanetsGetResponse "Список багажей"
// @Failure 400 {object} model.PlanetsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.PlanetsGetResponse "Внутренняя ошибка сервера"
// @Router /planets/{planet_id}/constellation [post]
func (h *Handler) RemovePlanetFromConstellation(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	searchCode := c.DefaultQuery("searchCode", "")
	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
		return
	}
	err = h.UseCase.RemovePlanetFromConstellation(uint(planetID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	planets, err := h.UseCase.GetPlanets(searchCode, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	constellation, err := h.UseCase.GetConstellationByIDUser(uint(planets.ConstellationID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"constellation": constellation})
}

// @Summary Добавление изображения к багажу
// @Description Добавляет изображение к багажу по его ID
// @Tags Планета
// @Accept mpfd
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Param image formData file true "Изображение планеты"
// @Success 200 {object} model.Planet "Информация о багаже с изображением"
// @Success 200 {object} model.Planet
// @Failure 400 {object} model.Planet "Некорректный запрос"
// @Failure 500 {object} model.Planet "Внутренняя ошибка сервера"
// @Router /planet/{planet_id}/image [post]
func (h *Handler) AddPlanetImage(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимое изображение"})
		return
	}

	file, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть изображение"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать изображение в байтах"})
		return
	}

	contentType := image.Header.Get("Content-Type")

	err = h.UseCase.AddPlanetImage(uint(planetID), uint(userID), imageBytes, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	planet, err := h.UseCase.GetPlanetByID(uint(planetID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": planet})
}
