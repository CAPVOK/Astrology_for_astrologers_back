package repository

import (
	"errors"
	"time"

	"space/internal/model"
)

func (r *Repository) GetConstellationsModerator(searchName, startFormationDate, endFormationDate, constellationStatus string, moderatorID uint) ([]model.ConstellationRequest, error) {
	query := r.db.Table("constellations").
		Select("DISTINCT constellations.constellation_id, constellations.constellation_name, constellations.start_date, constellations.end_date, constellations.creation_date, constellations.formation_date, constellations.confirmation_date, constellations.constellation_status, users.full_name").
		Joins("JOIN users ON users.user_id = constellations.user_id").
		Where("constellations.constellation_status LIKE ? AND constellations.constellation_name LIKE ? AND constellations.constellation_status != ?", constellationStatus, searchName, model.CONSTELLATION_STATUS_DELETED)
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
		Select("constellations.constellation_id, constellations.flight_number, constellations.creation_date, constellations.formation_date, constellations.confirmation_date, constellations.constellation_status, users.full_name").
		Joins("JOIN users ON users.user_id = constellations.user_id").
		Where("constellations.constellation_status != ? AND constellations.constellation_id = ?", model.CONSTELLATION_STATUS_DELETED, constellationID).
		Scan(&constellation).Error; err != nil {
		return model.ConstellationGetResponse{}, errors.New("ошибка получения созвездия по ИД")
	}
	var planets []model.Planet
	if err := r.db.
		Table("planets").
		Joins("JOIN constellation_planets ON planets.planet_id = constellation_planets.planet_id").
		Where("constellation_planets.constellation_id = ?", constellation.ConstellationID).
		Scan(&planets).Error; err != nil {
		return model.ConstellationGetResponse{}, errors.New("ошибка получения планет для созвездия")
	}
	constellation.Planets = planets
	return constellation, nil
}

func (r *Repository) UpdateConstellationModerator(constellationID uint, moderatorID uint, newConstellation model.ConstellationUpdateRequest) error {
	var constellation model.Constellation
	if err := r.db.Table("constellations").
		Where("constellation_id = ? AND moderator_id = ?", constellationID, moderatorID).
		First(&constellation).
		Error; err != nil {
		return errors.New("созвездие не найдена или не принадлежит указанному модератору")
	}
	if err := r.db.Table("constellations").
		Model(&constellation).
		Update("name", newConstellation.Name).
		Update("end_date", newConstellation.EndDate).
		Update("start_date", newConstellation.StartDate).
		Error; err != nil {
		return errors.New("ошибка обновления номера рейса")
	}
	return nil
}

func (r *Repository) UpdateConstellationStatusModerator(constellationID, moderatorID uint, constellationStatus model.ConstellationUpdateStatusRequest) error {
	var constellation model.Constellation
	if err := r.db.Table("constellations").
		Where("constellation_id = ? AND moderator_id = ? AND constellation_status = ?", constellationID, moderatorID, model.CONSTELLATION_STATUS_WORK).
		First(&constellation).
		Error; err != nil {
		return errors.New("созвездие не найдена или не принадлежит указанному модератору")
	}
	constellation.ConstellationStatus = constellationStatus.ConstellationStatus
	constellation.ConfirmationDate = time.Now()
	if err := r.db.Save(&constellation).Error; err != nil {
		return errors.New("ошибка обновления статуса созвездия в БД")
	}

	return nil
}
