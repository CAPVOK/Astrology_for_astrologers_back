package api

import (
	"net/http"
	"space/internal/app/ds"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPlanets(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Ошибка авторизации"})
		return
	}
	if isAdmin {
		searchQuery := c.DefaultQuery("searchByName", "")
		foundPlanets, err := h.Repo.SearchPlanetsByNameAdmin(searchQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось загрузить данные"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planets": foundPlanets})
	} else {
		searchQuery := c.DefaultQuery("searchByName", "")
		foundPlanets, err := h.Repo.SearchPlanetsByName(searchQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось загрузить данные"})
			return
		}
		constellation, err := h.Repo.GetCreatedConstellationByUser(USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "planets": foundPlanets, "constellationId": ""})
			return
		}
		newConst := map[string]interface{}{
			"Id": constellation.Id,
			/* "Name":             constellation.Name,
			"StartDate":        constellation.StartDate,
			"EndDate":          constellation.EndDate,
			"Status":           constellation.Status,
			"CreationDate":     constellation.CreationDate,
			"FormationDate":    constellation.FormationDate,
			"ConfirmationDate": constellation.ConfirmationDate, */
		}
		c.JSON(http.StatusOK, gin.H{"planets": foundPlanets, "constellationId": newConst})
	}
}

func (h *Handler) GetPlanetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Некорректные данные"})
		return
	}
	_, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Ошибка авторизации"})
		return
	}
	if isAdmin {
		planet, err := h.Repo.GetPlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Планета не найдена"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planet": planet})
	} else {
		planet, err := h.Repo.GetActivePlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Планета не найдена"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planet": planet})
	}

}

func (h *Handler) DeletePlanetById(c *gin.Context) {
	_, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Ошибка авторизации"})
		return
	}
	if isAdmin {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Некорректные данные"})
			return
		}
		_, err = h.Repo.GetPlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Планета не найдена"})
			return
		}
		err = h.Repo.DeactivatePlanetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось удалить планету"})
			return
		}
		planets, err := h.Repo.GetPlanets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось загрузить данные"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Планета удалена успешно", "planets": planets})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "Нет прав"})
		return
	}
}

func (h *Handler) ChangePlanetById(c *gin.Context) {
	_, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		var updatedPlanet ds.Planet
		if err := c.BindJSON(&updatedPlanet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Некорректные данные"})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Некорректные данные"})
			return
		}
		_, err = h.Repo.GetActivePlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Планета не найдена"})
			return
		}
		if err := h.Repo.UpdatePlanetByID(id, &updatedPlanet); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось обновить данные"})
			return
		}
		planet, err := h.Repo.GetActivePlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось загрузить данные"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Данные планеты обновлены успешно", "planet": planet})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "Нет прав"})
		return
	}
}

func (h *Handler) CreatePlanet(c *gin.Context) {
	_, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Ошибка авторизации"})
		return
	}
	if isAdmin {
		var newPlanet ds.Planet
		if err := c.BindJSON(&newPlanet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Некорректные данные"})
			return
		}
		if err := h.Repo.CreatePlanet(&newPlanet); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось создать планету"})
			return
		}
		planets, err := h.Repo.GetActivePlanets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось загрузить данные"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Планета создана успешно", "planets": planets})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "Нет прав"})
		return
	}
}

func (h *Handler) AddPlanetById(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Ошибка авторизации"})
		return
	}
	if isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "Нет прав"})
		return
	} else {
		planetId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Некорректные данные"})
			return
		}
		_, err = h.Repo.GetActivePlanetById(planetId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Планета не найдена"})
			return
		}
		constellation, err := h.Repo.GetCreatedConstellationByUser(USERID)
		if err != nil {
			err = h.Repo.CreateConstellationForUser(USERID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось создать созвездие"})
				return
			}
			constellation, err = h.Repo.GetCreatedConstellationByUser(USERID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось создать созвездие"})
				return
			}
		}
		if err := h.Repo.AddPlanetToConstellation(uint(planetId), constellation.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не удалось добавить планету"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Планета добавлена успешно"})
	}
}
