package delivery

import (
	"net/http"
	"strconv"

	"space/internal/model"
	"space/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// GetConstellations godoc
// @Summary Получение списка созвездий
// @Description Возвращает список всех не удаленных созвездий
// @Tags Созвездие
// @Produce json
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param constellationStatus query string false "Статус созвездия" Format(email)
// @Success 200 {object} model.ConstellationRequest "Список созвездий"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /constellation [get]
func (h *Handler) GetConstellations(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	/* searchFlightNumber := c.DefaultQuery("searchFlightNumber", "") */
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	constellationStatus := c.DefaultQuery("constellationStatus", "")
	var constellations []model.ConstellationRequest
	var err error
	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		constellations, err = h.UseCase.GetConstellationsModerator("", startFormationDate, endFormationDate, constellationStatus, userID)
	} else {
		constellations, err = h.UseCase.GetConstellationsUser("", startFormationDate, endFormationDate, constellationStatus, userID)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"constellations": constellations})
}

// GetConstellationByID godoc
// @Summary Получение созвездия по идентификатору
// @Description Возвращает информацию о созвездии по её идентификатору
// @Tags Созвездие
// @Produce json
// @Param constellation_id path int true "Идентификатор созвездия"
// @Success 200 {object} model.ConstellationGetResponse "Информация о созвездии"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /constellation/{constellation_id} [get]
func (h *Handler) GetConstellationByID(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД созвездия"})
		return
	}
	var constellation model.ConstellationGetResponse
	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		constellation, err = h.UseCase.GetConstellationByIDModerator(uint(constellationID), userID)
	} else {
		constellation, err = h.UseCase.GetConstellationByIDUser(uint(constellationID), userID)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if constellation.ConstellationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Созвездие не найдено"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"constellation": constellation})
}

// DeleteConstellation godoc
// @Summary Удаление созвездия
// @Description Удаляет доставку по её идентификатору
// @Tags Созвездие
// @Produce json
// @Param id path int true "Идентификатор созвездия"
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param constellationStatus query string false "Статус созвездия" Format(email)
// @Success 200 {object} model.ConstellationRequest "Список багажей"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /constellation/{constellation_id} [delete]
func (h *Handler) DeleteConstellation(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	/* searchFlightNumber := c.DefaultQuery("searchFlightNumber", "") */
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	constellationStatus := c.DefaultQuery("constellationStatus", "")
	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД созвездия"})
		return
	}
	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		err = h.UseCase.DeleteConstellationUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		constellations, err := h.UseCase.GetConstellationsModerator("", startFormationDate, endFormationDate, constellationStatus, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellations": constellations})
	} else {
		err = h.UseCase.DeleteConstellationUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		constellations, err := h.UseCase.GetConstellationsUser("", startFormationDate, endFormationDate, constellationStatus, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellations": constellations})
	}
}

// UpdateConstellation godoc
// @Summary Обновление созвездие
// @Description Обновляет поля созвездия по её идентификатору
// @Tags Созвездие
// @Produce json
// @Param constellation_id path int true "Идентификатор созвездия"
// @Param  newConstellation body model.ConstellationUpdateRequest true "Новое созвездие"
// @Success 200 {object} model.ConstellationGetResponse "Информация о созвездии"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /constellation/{constellation_id}/update [put]
func (h *Handler) UpdateConstellation(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД созвездия"})
		return
	}
	var constellation model.ConstellationUpdateRequest
	if err := c.BindJSON(&constellation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения JSON объекта"})
		return
	}
	err = h.UseCase.UpdateConstellationUser(uint(constellationID), userID, constellation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newConstellation, err := h.UseCase.GetConstellationByIDUser(uint(constellationID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"constellation": newConstellation})
	/* if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		err = h.UseCase.UpdateFlightNumberModerator(uint(constellationID), userID, constellation)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		constellation, err := h.UseCase.GetConstellationByIDModerator(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	} else {
		err = h.UseCase.UpdateConstellationUser(uint(constellationID), userID, flightNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		constellation, err := h.UseCase.GetConstellationByIDUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	} */
}

/* func (h *Handler) UpdateConstellationStatusUser(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недоупстимый ИД созвездия"})
		return
	}
	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос доступен только пользователю"})
		return
	} else {
		err = h.UseCase.UpdateConstellationStatusUser(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		constellation, err := h.UseCase.GetConstellationByIDUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	}
} */

/* func (h *Handler) UpdateConstellationStatusModerator(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД созвездия"})
		return
	}

	var constellationStatus model.ConstellationUpdateStatusRequest
	if err := c.BindJSON(&constellationStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		err = h.UseCase.UpdateConstellationStatusModerator(uint(constellationID), userID, constellationStatus)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		constellation, err := h.UseCase.GetConstellationByIDUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос доступен только модератору"})
		return
	}
} */

// UpdateConstellationStatus godoc
// @Summary Обновление статуса созвездия
// @Description Обновляет статус чернового созвездия для юзера и обновляет статус созвездия по идентификатору созвездия и новому статусу для модератора
// @Tags Созвездие
// @Produce json
// @Param constellation_id path int true "Идентификатор созвездия"
// @Param constellationStatus body model.ConstellationUpdateStatusRequest true "Новый статус созвездия"
// @Success 200 {object} model.ConstellationGetResponse "Информация о созвездии"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /constellation/{constellation_id}/status [put]
func (h *Handler) UpdateConstellationStatus(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД созвездия"})
		return
	}
	var constellationStatus model.ConstellationUpdateStatusRequest
	if err := c.BindJSON(&constellationStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// модератор
	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		err = h.UseCase.UpdateConstellationStatusModerator(uint(constellationID), userID, constellationStatus)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		constellation, err := h.UseCase.GetConstellationByIDModerator(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if constellation.ConstellationID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Созвездие не найдено"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	} else {
		if constellationStatus.ConstellationStatus == model.CONSTELLATION_STATUS_REJECTED || constellationStatus.ConstellationStatus == model.CONSTELLATION_STATUS_COMPLETED {
			c.JSON(http.StatusForbidden, gin.H{"error": "не хватает прав"})
			return
		}
		// юзер
		id, err := h.UseCase.UpdateConstellationStatusUser(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		constellation, err := h.UseCase.GetConstellationByIDUser(id, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	}
}
