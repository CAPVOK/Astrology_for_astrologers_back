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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
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
		constellation, err := h.Repo.GetCreatedConstellationByUser(USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "planets": foundPlanets, "constellationId": "", "message": "Constellation not found"})
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
		c.JSON(http.StatusOK, gin.H{"planets": foundPlanets, "constellation": newConst})
	}
}

func (h *Handler) GetPlanetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		planet, err := h.Repo.GetPlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Planet not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planet": planet})
	} else {
		planet, err := h.Repo.GetActivePlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Planet not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"planet": planet})
	}

}

func (h *Handler) DeletePlanetById(c *gin.Context) {
	_, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = h.Repo.GetPlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Planet not found"})
			return
		}
		err = h.Repo.DeactivatePlanetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to delete the planet"})
			return
		}
		planets, err := h.Repo.GetPlanets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Planet deleted successfully", "planets": planets})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = h.Repo.GetActivePlanetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Planet not found"})
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
		c.JSON(http.StatusOK, gin.H{"message": "Planet updated successfully", "planet": planet})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	}
}

func (h *Handler) CreatePlanet(c *gin.Context) {
	_, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		var newPlanet ds.Planet
		if err := c.BindJSON(&newPlanet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
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
		c.JSON(http.StatusCreated, gin.H{"message": "Planet created successfully", "planets": planets})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	}
}

func (h *Handler) AddPlanetById(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	} else {
		planetId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = h.Repo.GetActivePlanetById(planetId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Planet not found"})
			return
		}
		constellation, err := h.Repo.GetCreatedConstellationByUser(USERID)
		if err != nil {
			err = h.Repo.CreateConstellationForUser(USERID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Fail creation constellation"})
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
		c.JSON(http.StatusOK, gin.H{"message": "Planet added to constellation successfully"})
	}
}
