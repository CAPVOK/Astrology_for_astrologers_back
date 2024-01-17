package repository

import (
	"errors"
	"time"

	"space/internal/model"
)

func (r *Repository) GetConstellationsModerator(searchName, startFormationDate, endFormationDate, constellationStatus string, moderatorID uint) ([]model.ConstellationRequest, error) {
	query := r.db.Table("constellations").
		Select("DISTINCT constellations.constellation_id, constellations.name, constellations.start_date, constellations.end_date, constellations.start_date, constellations.end_date, constellations.creation_date, constellations.formation_date, constellations.confirmation_date, constellations.constellation_status, users.full_name").
		Joins("JOIN users ON users.user_id = constellations.user_id").
		Where("constellations.constellation_status LIKE ? AND constellations.name LIKE ? AND constellations.constellation_status != ? AND constellations.constellation_status != ?", constellationStatus, searchName, model.CONSTELLATION_STATUS_DELETED, model.CONSTELLATION_STATUS_DRAFT)
	if startFormationDate != "" && endFormationDate != "" {
		query = query.Where("constellations.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
	}
	var constellations []model.ConstellationRequest
	if err := query.Find(&constellations).Error; err != nil {
		return nil, errors.New("ошибка получения cозвездий")
	}
	return constellations, nil
}

func (r *Repository) GetConstellationByIDModerator(constellationID, moderatorID uint) (model.ConstellationGetResponse, error) {
	var constellation model.ConstellationGetResponse
	if err := r.db.
		Table("constellations").
		Select("constellations.constellation_id, constellations.name, constellations.start_date, constellations.end_date, constellations.creation_date, constellations.formation_date, constellations.confirmation_date, constellations.constellation_status, users.full_name").
		Joins("JOIN users ON users.user_id = constellations.user_id").
		Where("constellations.constellation_status != ? AND constellations.constellation_id = ?", model.CONSTELLATION_STATUS_DELETED, constellationID).
		Scan(&constellation).Error; err != nil {
		return model.ConstellationGetResponse{}, errors.New("не удалось найти созвездие. Ошибка 831781")
	}
	var planets []model.PlanetInConstellation
	if err := r.db.
		Table("planets").
		Joins("JOIN constellation_planets ON planets.planet_id = constellation_planets.planet_id").
		Where("constellation_planets.constellation_id = ?", constellationID).
		Scan(&planets).Error; err != nil {
		return model.ConstellationGetResponse{}, errors.New("ошибка получения планет для созвездия")
	}
	constellation.Planets = planets
	if constellation.ConstellationID == 0 {
		return constellation, errors.New("не удалось найти созвездие. Ошибка 882820")
	}
	return constellation, nil
}

func (r *Repository) UpdateConstellationStatusModerator(constellationID, moderatorID uint, constellationStatus model.ConstellationUpdateStatusRequest) error {
	var constellation model.Constellation
	if err := r.db.Table("constellations").
		Where("constellation_id = ? AND constellation_status = ?", constellationID, model.CONSTELLATION_STATUS_WORK).
		First(&constellation).
		Error; err != nil {
		return errors.New("созвездие не найдено")
	}
	constellation.ConstellationStatus = constellationStatus.ConstellationStatus
	constellation.ModeratorID = &moderatorID
	currentTime := time.Now()
	constellation.ConfirmationDate = &currentTime
	if err := r.db.Save(&constellation).Error; err != nil {
		return errors.New("ошибка обновления статуса созвездия в БД")
	}

	return nil
}
