package delivery

import (
	"io"
	"net/http"
	"strconv"

	"space/internal/model"
	"space/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка планет
// @Description Возращает список всех активных планет и ид черновой заявки
// @Tags Планета
// @Produce json
// @Param searchName query string false "Название планеты" Format(email)
// @Success 200 {object} model.PlanetsGetResponse "Список планет"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet [get]
func (h *Handler) GetPlanets(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте пп"})
		return
	}
	userID := ctxUserID.(uint)
	searchCode := c.DefaultQuery("searchName", "")

	planets, err := h.UseCase.GetPlanets(searchCode, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planets": planets.Planets, "constellationID": planets.ConstellationID})
}

// @Summary Получение планеты по ID
// @Description Возвращает информацию о планете по ее ID
// @Tags Планета
// @Produce json
// @Param planet_id path int true "ID планетаа"
// @Success 200 {object} model.Planet "Информация о планетае"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{planet_id} [get]
func (h *Handler) GetPlanetByID(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": planet})
}

// @Summary Создание новой планеты
// @Description Создает новую планету с предоставленными данными
// @Tags Планета
// @Accept json
// @Produce json
// @Param searchName query string false "Название планеты" Format(email)
// @Param planet body model.PlanetRequest true "Пользовательский объект в формате JSON"
// @Success 200 {object} model.PlanetsGetResponse "Список планет"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} model.ErrorResponse "У пользователя нет прав для этого запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet [post]
func (h *Handler) CreatePlanet(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	searchCode := c.DefaultQuery("searchName", "")

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		var planet model.PlanetRequest
		if err := c.BindJSON(&planet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
			return
		}
		err := h.UseCase.CreatePlanet(userID, planet)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		planets, err := h.UseCase.GetPlanets(searchCode, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planets": planets.Planets, "constellationID": planets.ConstellationID})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет этих прав"})
		return
	}
}

// @Summary Удаление планеты
// @Description Удаляет планету по его ID
// @Tags Планета
// @Produce json
// @Param planet_id path int true "ID планетаа"
// @Param searchName query string false "Название планеты" Format(email)
// @Success 200 {object} model.PlanetsGetResponse "Список планет"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} model.ErrorResponse "У пользователя нет прав для этого запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{planet_id} [delete]
func (h *Handler) DeletePlanet(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	searchCode := c.DefaultQuery("searchName", "")

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		planetID, err := strconv.Atoi(c.Param("planet_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
			return
		}
		err = h.UseCase.DeletePlanet(uint(planetID), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		planets, err := h.UseCase.GetPlanets(searchCode, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planets": planets.Planets, "constellationID": planets.ConstellationID})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет этих прав"})
		return
	}

}

// @Summary Обновление информации о планетe
// @Description Обновляет информацию о планетe по его ID
// @Tags Планета
// @Accept json
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Success 200 {object} model.Planet "Информация о планетe"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} model.ErrorResponse "У пользователя нет прав для этого запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{planet_id} [put]
func (h *Handler) UpdatePlanet(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "недопустимый ИД планеты"}})
		return
	}
	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		var planet model.PlanetRequest
		if err := c.BindJSON(&planet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
			return
		}

		err = h.UseCase.UpdatePlanet(uint(planetID), uint(userID), planet)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updatedPlanet, err := h.UseCase.GetPlanetByID(uint(planetID), uint(userID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"planet": updatedPlanet})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет этих прав"})
		return
	}

}

// @Summary Добавление планеты к созвездии
// @Description Добавляет планету к созвездию по ее ID
// @Tags Планета
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Param searchName query string false "Код планеты" Format(email)
// @Success 200 {object} model.PlanetsGetResponse  "Список планет"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{planet_id}/constellation	 [post]
func (h *Handler) AddPlanetToConstellation(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	searchCode := c.DefaultQuery("searchName", "")

	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
		return
	}
	err = h.UseCase.AddPlanetToConstellation(uint(planetID), uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	planets, err := h.UseCase.GetPlanets(searchCode, uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"planets": planets.Planets, "constellationID": planets.ConstellationID})
}

// @Summary Удаление планеты из созвездия
// @Description Удаляет планета из созвездия по еe ID
// @Tags Планета
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Param searchName query string false "Код планеты" Format(email)
// @Success 200 {object} model.PlanetsGetResponse "Список планет"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{planet_id}/constellation [delete]
func (h *Handler) RemovePlanetFromConstellation(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	searchCode := c.DefaultQuery("searchName", "")
	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
		return
	}
	err = h.UseCase.RemovePlanetFromConstellation(uint(planetID), uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	planets, err := h.UseCase.GetPlanets(searchCode, uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	constellation, err := h.UseCase.GetConstellationByIDUser(uint(planets.ConstellationID), uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"constellation": constellation})
}

// @Summary Добавление изображения к планете
// @Description Добавляет изображение к планете по ее ID
// @Tags Планета
// @Accept mpfd
// @Produce json
// @Param planet_id path int true "ID планета"
// @Param image formData file true "Изображение планеты"
// @Success 200 {object} model.Planet "Информация о планете с изображением"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} model.ErrorResponse "У пользователя нет прав для этого запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{planet_id}/image [post]
func (h *Handler) AddPlanetImage(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	planetID, err := strconv.Atoi(c.Param("planet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД планеты"})
		return
	}
	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		image, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимое изображение"})
			return
		}
		file, err := image.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось открыть изображение"})
			return
		}
		defer file.Close()
		imageBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать изображение в байтах"})
			return
		}
		contentType := image.Header.Get("Content-Type")
		err = h.UseCase.AddPlanetImage(uint(planetID), uint(userID), imageBytes, contentType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		planet, err := h.UseCase.GetPlanetByID(uint(planetID), uint(userID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planet": planet})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет этих прав"})
		return
	}
}
