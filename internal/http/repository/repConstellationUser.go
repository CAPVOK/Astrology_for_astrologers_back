package repository

import (
	"errors"
	"time"

	"space/internal/model"
)

func (r *Repository) GetConstellationsUser(searchFlightNumber, startFormationDate, endFormationDate, constellationStatus string, userID uint) ([]model.ConstellationRequest, error) {
	query := r.db.Table("constellations").
		Select("DISTINCT constellations.constellation_id, constellations.flight_number, constellations.creation_date, constellations.formation_date, constellations.confirmation_date, constellations.constellation_status, users.full_name").
		Joins("JOIN users ON users.user_id = constellations.user_id").
		Where("constellations.constellation_status LIKE ? AND constellations.flight_number LIKE ? AND constellations.user_id = ? AND constellations.constellation_status != ?", constellationStatus, searchFlightNumber, userID, model.CONSTELLATION_STATUS_DELETED)

	if startFormationDate != "" && endFormationDate != "" {
		query = query.Where("constellations.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
	}

	var constellations []model.ConstellationRequest
	if err := query.Find(&constellations).Error; err != nil {
		return nil, errors.New("ошибка получения созвездий")
	}

	return constellations, nil
}

func (r *Repository) GetConstellationByIDUser(constellationID, userID uint) (model.ConstellationGetResponse, error) {
	var constellation model.ConstellationGetResponse

	if err := r.db.
		Table("constellations").
		Select("constellations.constellation_id, constellations.flight_number, constellations.creation_date, constellations.formation_date, constellations.confirmation_date, constellations.constellation_status, users.full_name").
		Joins("JOIN users ON users.user_id = constellations.user_id").
		Where("constellations.constellation_status != ? AND constellations.constellation_id = ? AND constellations.user_id = ?", model.CONSTELLATION_STATUS_DELETED, constellationID, userID).
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

func (r *Repository) DeleteConstellationUser(constellationID, userID uint) error {
	var constellation model.Constellation
	if err := r.db.Table("constellations").
		Where("constellation_id = ? AND user_id = ?", constellationID, userID).
		First(&constellation).
		Error; err != nil {
		return errors.New("созвездие не найдена или не принадлежит указанному пользователю")
	}

	tx := r.db.Begin()
	if err := tx.Where("constellation_id = ?", constellationID).Delete(&model.ConstellationPlanet{}).Error; err != nil {
		tx.Rollback()
		return errors.New("ошибка удаления связей из таблицы-множества")
	}

	err := r.db.Model(&model.Constellation{}).Where("constellation_id = ?", constellationID).Update("constellation_status", model.CONSTELLATION_STATUS_DELETED).Error
	if err != nil {
		return errors.New("ошибка обновления статуса на удален")
	}
	tx.Commit()

	return nil
}

func (r *Repository) UpdateConstellationUser(constellationID uint, userID uint, newConstellation model.ConstellationUpdateRequest) error {
	var constellation model.Constellation
	if err := r.db.Table("constellations").
		Where("constellation_id = ? AND user_id = ?", constellationID, userID).
		First(&constellation).
		Error; err != nil {
		return errors.New("созвездие не найдена или не принадлежит указанному пользователю")
	}

	if err := r.db.Table("constellations").
		Model(&constellation).
		Update("name", newConstellation.Name).
		Update("end_date", newConstellation.EndDate).
		Update("start_date", newConstellation.StartDate).
		Error; err != nil {
		return errors.New("ошибка обновления созвездия")
	}

	return nil
}

func (r *Repository) UpdateConstellationStatusUser(constellationID, userID uint) error {
	var constellation model.Constellation
	if err := r.db.Table("constellations").
		Where("constellation_id = ? AND user_id = ? AND constellation_status = ?", constellationID, userID, model.CONSTELLATION_STATUS_DRAFT).
		First(&constellation).
		Error; err != nil {
		return errors.New("созвездие не найдена, или не принадлежит указанному пользователю, или не имеет статус черновик")
	}

	constellation.ConstellationStatus = model.CONSTELLATION_STATUS_WORK
	constellation.FormationDate = time.Now()

	if err := r.db.Save(&constellation).Error; err != nil {
		return errors.New("ошибка обновления статуса созвездия в БД")
	}

	return nil
}
