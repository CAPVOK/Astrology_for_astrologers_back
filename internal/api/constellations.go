package api

import (
	"net/http"
	"space/internal/app/ds"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetConstellations(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		constellations, err := h.Repo.GetActiveConstellations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellations": constellations})
	} else {
		constellations, err := h.Repo.GetActiveConstellationsByUser(USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellations": constellations})
	}
}

func (h *Handler) GetConstellationById(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constellation ID"})
		return
	}
	if isAdmin {
		constellation, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	} else {
		constellation, err := h.Repo.GetConstellationById(id, USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"constellation": constellation})
	}
}

func (h *Handler) DeleteConstellationById(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constellation ID"})
		return
	}
	if isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	} else {
		_, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := h.Repo.UpdateStatusToDeleted(id, USERID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := h.Repo.DeleteAllPlanetsFromConstellation(id, USERID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		constellations, err := h.Repo.GetActiveConstellationsByUser(USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Status updated to 'deleted'", "constellations": constellations})
	}
}

func (h *Handler) ChangeConstellationById(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constellation ID"})
		return
	}
	if isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	} else {
		var updatedConstellation ds.Constellation
		if err := c.BindJSON(&updatedConstellation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}
		constellation, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if constellation.Status != "created" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update constellation that is not in 'created' status"})
			return
		}
		if err := h.Repo.UpdateConstellationByID(id, &updatedConstellation); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		newConstellation, err := h.Repo.GetConstellationById(id, USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Сonstellations updated successfully", "constellation": newConstellation})
	}
}

func (h *Handler) DoConstellationInProgress(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	} else {
		constellation, err := h.Repo.GetCreatedConstellationByUser(USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "There is no created constellation"})
			return
		}
		if constellation.Status != "created" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update status to 'inprogress' for a constellation that is not in 'created' status"})
			return
		}
		if err := h.Repo.UpdateStatusToInProgress(int(constellation.Id), uint(USERID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed making status 'inprogress'"})
			return
		}
		newConstellation, err := h.Repo.GetConstellationById(int(constellation.Id), USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Status updated to 'inprogress'", "constellation": newConstellation})
	}
}

func (h *Handler) DoConstelltionCanceledById(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constellation ID"})
			return
		}
		constellation, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if constellation.Status != "inprogress" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update status to 'canceled' for a constellation that is not in 'inprogress' status"})
			return
		}
		if err := h.Repo.UpdateStatusToCanceled(id, uint(USERID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		constellation, err = h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Status updated to 'canceled'", "constellation": constellation})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	}
}

func (h *Handler) DoConstelltionCompletedById(c *gin.Context) {
	USERID, isAdmin, err := singleton()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "auth error"})
		return
	}
	if isAdmin {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constellation ID"})
			return
		}
		constellation, err := h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if constellation.Status != "inprogress" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update status to 'completed' for a constellation that is not in 'inprogress' status"})
			return
		}
		if err := h.Repo.UpdateStatusToCompleted(id, uint(USERID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		constellation, err = h.Repo.GetConstellationByIdAdmin(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Status updated to 'completed'", "constellation": constellation})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "no rules"})
		return
	}
}

func (h *Handler) RemovePlanetById(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planet ID"})
			return
		}
		constellation, err := h.Repo.GetCreatedConstellationByUser(USERID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "no created constellatin"})
			return
		}
		curConstellation, err := h.Repo.GetConstellationByIdAdmin(int(constellation.Id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "no created constellation"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "There is no planet with such id in the created constellation"})
			return
		}
		if err := h.Repo.DeletePlanetFromConstellation(uint(planetId), constellation.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Planet deleted successfully"})
	}
}
