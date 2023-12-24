package api

import (
	"net/http"
	"space/internal/app/ds"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetConstellations godoc
// @Summary      returns the list of constellations
// @Description  returns the list of constellations
// @Tags         expeditions
// @Produce      json
// @Success      200  {object} object{constellations=ds.Constellation}
// @Failure      400  {object} object{status=string, message=string}
// @Failure      500  {object} object{status=string, message=string}
// @Router       /constellation/ [get]
func (h *Handler) GetConstellations(c *gin.Context) {
	value, exists := c.Get("sessionContext")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Вы должны быть авторизованы",
		})
		return
	}
	sc := value.(ds.SessionContext)

	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	status := c.DefaultQuery("status", "")

	if sc.Role == ds.Moderator {
		constellations, err := h.Repo.GetActiveConstellations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellations": constellations})
	} else {
		constellations, err := h.Repo.GetActiveConstellationsByUser(sc.UserID, startFormationDate, endFormationDate, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellations": constellations})
	}
}

func (h *Handler) GetConstellationById(c *gin.Context) {
	value, exists := c.Get("sessionContext")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Вы должны быть авторизованы",
		})
		return
	}
	sc := value.(ds.SessionContext)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id созвездия"})
		return
	}
	if sc.Role == ds.Moderator {
		constellation, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	} else {
		constellation, err := h.Repo.GetConstellationById(id, sc.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	}
}

func (h *Handler) DeleteConstellationById(c *gin.Context) {
	value, exists := c.Get("sessionContext")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "вы должны быть авторизованы",
		})
		return
	}
	sc := value.(ds.SessionContext)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id созвездия"})
		return
	}
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	status := c.DefaultQuery("status", "")
	if sc.Role == ds.Moderator {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	} else {
		_, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := h.Repo.UpdateStatusToDeleted(id, sc.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := h.Repo.DeleteAllPlanetsFromConstellation(id, sc.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		constellations, err := h.Repo.GetActiveConstellationsByUser(sc.UserID, startFormationDate, endFormationDate, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Удалено", "constellations": constellations})
	}
}

func (h *Handler) ChangeConstellationById(c *gin.Context) {
	value, exists := c.Get("sessionContext")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Необходимо авторизоваться",
		})
		return
	}
	sc := value.(ds.SessionContext)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id созвездия"})
		return
	}
	if sc.Role == ds.Moderator {
		c.JSON(http.StatusForbidden, gin.H{"message": "недостаточно прав"})
		return
	} else {
		var updatedConstellation ds.Constellation
		if err := c.BindJSON(&updatedConstellation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при получении данных"})
			return
		}
		constellation, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if constellation.Status != "created" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при обновлении статуса"})
			return
		}
		if err := h.Repo.UpdateConstellationByID(id, &updatedConstellation); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		newConstellation, err := h.Repo.GetConstellationById(id, sc.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "созвездие обновлено успешно", "constellation": newConstellation})
	}
}

func (h *Handler) DoConstellationInProgress(c *gin.Context) {
	value, exists := c.Get("sessionContext")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "необходимо авторизоваться",
		})
		return
	}
	sc := value.(ds.SessionContext)
	if sc.Role == ds.Moderator {
		c.JSON(http.StatusForbidden, gin.H{"message": "недостатчно прав"})
		return
	} else {
		constellation, err := h.Repo.GetCreatedConstellationByUser(sc.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "нету только созданных созвездий"})
			return
		}
		if constellation.Status != "created" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при обновлении статуса"})
			return
		}
		if err := h.Repo.UpdateStatusToInProgress(int(constellation.Id), uint(sc.UserID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "ошибка при обновлении статуса"})
			return
		}
		newConstellation, err := h.Repo.GetConstellationById(int(constellation.Id), sc.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "статус обновлен", "constellation": newConstellation})
	}
}

// DoConstelltionCanceledById godoc
// @Summary      cancels the cancellation
// @Description  cancels the cancellation by id
// @Tags         expeditions
// @Produce      json
// @Param        id path uint true "id of constellation"
// @Success      200  {object} object{message=string,constellations=ds.Constellation}
// @Failure      400  {object} object{status=string,message=string}
// @Failure      500  {object} object{status=string,message=string}
// @Router       /constellation/cancel/{id} [get]
func (h *Handler) DoConstelltionCanceledById(c *gin.Context) {
	value, exists := c.Get("sessionContext")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "необходимо авторизоваться",
		})
		return
	}
	sc := value.(ds.SessionContext)
	if sc.Role == ds.Moderator {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id созвездия"})
			return
		}
		constellation, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if constellation.Status != "inprogress" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка обновления статуса"})
			return
		}
		if err := h.Repo.UpdateStatusToCanceled(id, uint(sc.UserID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		constellation, err = h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "статус обновлен", "constellation": constellation})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "недостаточно прав"})
		return
	}
}

func (h *Handler) DoConstelltionCompletedById(c *gin.Context) {
	value, exists := c.Get("sessionContext")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "необходимо авторизоваться",
		})
		return
	}
	sc := value.(ds.SessionContext)
	if sc.Role == ds.Moderator {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверное id созвездия"})
			return
		}
		constellation, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if constellation.Status != "inprogress" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка обновления статуса"})
			return
		}
		if err := h.Repo.UpdateStatusToCompleted(id, uint(sc.UserID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		constellation, err = h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "статус обновлен", "constellation": constellation})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	}
}

func (h *Handler) RemovePlanetById(c *gin.Context) {
	value, exists := c.Get("sessionContext")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "необходимо авторизоваться",
		})
		return
	}
	sc := value.(ds.SessionContext)
	if sc.Role == ds.Moderator {
		c.JSON(http.StatusForbidden, gin.H{"message": "недостаточно прав"})
		return
	} else {
		planetId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверное id планеты"})
			return
		}
		constellation, err := h.Repo.GetCreatedConstellationByUser(sc.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "нету созданных созвездий"})
			return
		}
		curConstellation, err := h.Repo.GetConstellationByIdAdmin(int(constellation.Id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "нету созданных созвездий"})
			return
		}
		isPlanetInConst := false
		for _, planet := range curConstellation.Planets {
			if planet.PlanetID == uint(planetId) {
				isPlanetInConst = true
				break
			}
		}
		if !isPlanetInConst {
			c.JSON(http.StatusBadRequest, gin.H{"error": "планеты с заданным id не существует"})
			return
		}
		if err := h.Repo.DeletePlanetFromConstellation(uint(planetId), constellation.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "планета удалена"})
	}
}
