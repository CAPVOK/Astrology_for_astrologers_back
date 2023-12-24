package api

import (
	"net/http"
	"space/internal/app/ds"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPlanets(c *gin.Context) {
	USERID := c.GetInt("userID")
	isAdmin := c.GetInt("role")
	if isAdmin == 2 {
		searchQuery := c.DefaultQuery("q", "")
		foundPlanets, err := h.Repo.SearchPlanetsByNameAdmin(searchQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planets": foundPlanets})
	} else {
		searchQuery := c.DefaultQuery("q", "")
		foundPlanets, err := h.Repo.SearchPlanetsByName(searchQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var constellation *ds.Constellation
		if USERID != 0 {
			constellation, err = h.Repo.GetCreatedConstellationByUser(USERID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "planets": foundPlanets, "constellationId": "", "message": "Созвездие не найдено"})
				return
			}
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
		c.JSON(http.StatusOK, gin.H{"planets": foundPlanets, "constellation": newConst})
	}
}

func (h *Handler) GetPlanetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isAdmin := c.GetInt("role")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Не авторизован"})
		return
	}
	if isAdmin == 2 {
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
	isAdmin := c.GetInt("role")
	if isAdmin == 2 {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = h.Repo.GetPlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Планета не найдена"})
			return
		}
		err = h.Repo.DeactivatePlanetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Ошибка при удалении планеты"})
			return
		}
		planets, err := h.Repo.GetPlanets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Планета удалена", "planets": planets})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	}
}

func (h *Handler) ChangePlanetById(c *gin.Context) {
	isAdmin := c.GetInt("role")

	if isAdmin == 2 {
		var updatedPlanet ds.Planet
		if err := c.BindJSON(&updatedPlanet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка запроса"})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = h.Repo.GetActivePlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Планета не найдена"})
			return
		}
		if err := h.Repo.UpdatePlanetByID(id, &updatedPlanet); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		planet, err := h.Repo.GetActivePlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Планета обновлена успешно", "planet": planet})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "Недостаточно прав"})
		return
	}
}

func (h *Handler) CreatePlanet(c *gin.Context) {
	isAdmin := c.GetInt("role")
	if isAdmin == 2 {
		var newPlanet ds.Planet
		if err := c.BindJSON(&newPlanet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка запроса"})
			return
		}
		if err := h.Repo.CreatePlanet(&newPlanet); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		planets, err := h.Repo.GetActivePlanets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Планета создана успешно", "planets": planets})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	}
}

func (h *Handler) AddPlanetById(c *gin.Context) {
	USERID := c.GetInt("userID")
	isAdmin := c.GetInt("role")
	if isAdmin == 2 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Недостаточно прав"})
		return
	} else {
		planetId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Ошибка при создании созвездия"})
				return
			}
			constellation, err = h.Repo.GetCreatedConstellationByUser(USERID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if err := h.Repo.AddPlanetToConstellation(uint(planetId), constellation.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Планета добавлена в созвездие успешно"})
	}
}
