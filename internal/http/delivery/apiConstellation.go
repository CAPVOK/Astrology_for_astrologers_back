package delivery

import (
	"net/http"
	"strconv"

	"space/internal/model"
	"space/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// GetConstellations godoc
// @Summary Получение списка доставок
// @Description Возвращает список всех не удаленных доставок
// @Tags Созвездие
// @Produce json
// @Param searchFlightNumber query string false "Номер рейса" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param constellationStatus query string false "Статус созвездия" Format(email)
// @Success 200 {object} model.ConstellationRequest "Список доставок"
// @Failure 500 {object} model.ConstellationRequest "Ошибка сервера"
// @Router /constellation [get]
func (h *Handler) GetConstellations(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	constellationStatus := c.DefaultQuery("constellationStatus", "")

	var constellations []model.ConstellationRequest
	var err error

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		constellations, err = h.UseCase.GetConstellationsModerator(searchFlightNumber, startFormationDate, endFormationDate, constellationStatus, userID)
	} else {
		constellations, err = h.UseCase.GetConstellationsUser(searchFlightNumber, startFormationDate, endFormationDate, constellationStatus, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"constellations": constellations})
}

// GetConstellationByID godoc
// @Summary Получение созвездия по идентификатору
// @Description Возвращает информацию о доставке по её идентификатору
// @Tags Созвездие
// @Produce json
// @Param id path int true "Идентификатор созвездия"
// @Success 200 {object} model.ConstellationGetResponse "Информация о доставке"
// @Failure 400 {object} model.ConstellationGetResponse "Недопустимый идентификатор созвездия"
// @Failure 500 {object} model.ConstellationGetResponse "Ошибка сервера"
// @Router /constellation/{id} [get]
func (h *Handler) GetConstellationByID(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, constellation)
}

// DeleteConstellation godoc
// @Summary Удаление созвездия
// @Description Удаляет доставку по её идентификатору
// @Tags Созвездие
// @Produce json
// @Param id path int true "Идентификатор созвездия"
// @Param searchFlightNumber query string false "Номер рейса" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param constellationStatus query string false "Статус созвездия" Format(email)
// @Success 200 {object} model.ConstellationRequest "Список багажей"
// @Failure 400 {object} model.ConstellationRequest "Недопустимый идентификатор созвездия"
// @Failure 500 {object} model.ConstellationRequest "Ошибка сервера"
// @Router /constellation/{id}/delete [delete]
func (h *Handler) DeleteConstellation(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	constellationStatus := c.DefaultQuery("constellationStatus", "")
	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД созвездия"})
		return
	}

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос недоступен для модератора"})
		return
	}

	err = h.UseCase.DeleteConstellationUser(uint(constellationID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	constellations, err := h.UseCase.GetConstellationsUser(searchFlightNumber, startFormationDate, endFormationDate, constellationStatus, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"constellations": constellations})
}

// UpdateConstellationFlightNumber godoc
// @Summary Обновление номера рейса созвездия
// @Description Обновляет номер рейса для созвездия по её идентификатору
// @Tags Созвездие
// @Produce json
// @Param id path int true "Идентификатор созвездия"
// @Param flightNumber body model.ConstellationUpdateFlightNumberRequest true "Новый номер рейса"
// @Success 200 {object} model.ConstellationGetResponse "Информация о доставке"
// @Failure 400 {object} model.ConstellationGetResponse "Недопустимый идентификатор созвездия или ошибка чтения JSON объекта"
// @Failure 500 {object} model.ConstellationGetResponse "Ошибка сервера"
// @Router /constellation/{id}/update [put]
func (h *Handler) UpdateConstellationFlightNumber(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД созвездия"})
		return
	}

	var flightNumber model.ConstellationUpdateRequest
	if err := c.BindJSON(&flightNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения JSON объекта"})
		return
	}

	if middleware.ModeratorOnly(h.UseCase.Repository, c) {
		err = h.UseCase.UpdateFlightNumberModerator(uint(constellationID), userID, flightNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		constellation, err := h.UseCase.GetConstellationByIDModerator(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	} else {
		err = h.UseCase.UpdateConstellationUser(uint(constellationID), userID, flightNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		constellation, err := h.UseCase.GetConstellationByIDUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	}
}

// UpdateConstellationStatusUser godoc
// @Summary Обновление статуса созвездия для пользователя
// @Description Обновляет статус созвездия для пользователя по идентификатору созвездия
// @Tags Созвездие
// @Produce json
// @Param id path int true "Идентификатор созвездия"
// @Success 200 {object} model.ConstellationGetResponse "Информация о доставке"
// @Failure 400 {object} model.ConstellationGetResponse "Недопустимый идентификатор созвездия"
// @Failure 500 {object} model.ConstellationGetResponse "Ошибка сервера"
// @Router /constellation/{id}/user [put]
func (h *Handler) UpdateConstellationStatusUser(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
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
		err = h.UseCase.UpdateConstellationStatusUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		constellation, err := h.UseCase.GetConstellationByIDUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	}
}

// UpdateConstellationStatusModerator godoc
// @Summary Обновление статуса созвездия для модератора
// @Description Обновляет статус созвездия для модератора по идентификатору созвездия
// @Tags Созвездие
// @Produce json
// @Param id path int true "Идентификатор созвездия"
// @Param constellationStatus body model.ConstellationUpdateStatusRequest true "Новый статус созвездия"
// @Success 200 {object} model.ConstellationGetResponse "Информация о доставке"
// @Failure 400 {object} model.ConstellationGetResponse "Недопустимый идентификатор созвездия или ошибка чтения JSON объекта"
// @Failure 500 {object} model.ConstellationGetResponse "Ошибка сервера"
// @Router /constellation/{id}/status [put]
func (h *Handler) UpdateConstellationStatusModerator(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		constellation, err := h.UseCase.GetConstellationByIDUser(uint(constellationID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос доступен только модератору"})
		return
	}
}
