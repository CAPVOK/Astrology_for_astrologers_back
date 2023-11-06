package repository

import (
	"space/internal/app/ds"
	"time"
)

/* созвездия */

/* список всех активных заявок */
func (r *Repository) GetActiveConstellations() (*[]ds.Constellation, error) {
	var constellations []ds.Constellation
	/* var constellationsClient []ConstellationClient */

	if err := r.db.Where("status != 'deleted'").Order("creation_date").Find(&constellations).Error; err != nil {
		return nil, err
	}
	/* for _, constellation := range constellations {
		constellationClient := ConstellationClient{
			Id:               constellation.Id,
			Name:             constellation.Name,
			StartDate:        constellation.StartDate,
			EndDate:          constellation.EndDate,
			Status:           constellation.Status,
			CreationDate:     constellation.CreationDate,
			FormationDate:    constellation.FormationDate,
			ConfirmationDate: constellation.ConfirmationDate,
		}
		constellationsClient = append(constellationsClient, constellationClient)
	} */
	return &constellations, nil
}

/* список активных созвездий для юзера */
func (r *Repository) GetActiveConstellationsByUser(userId int) (*[]ConstellationClient, error) {
	var constellations []ds.Constellation
	var constellationsClient []ConstellationClient
	if err := r.db.Where("status != 'deleted' AND user_id = ?", userId).Order("creation_date").Find(&constellations).Error; err != nil {
		return nil, err
	}
	for _, constellation := range constellations {
		constellationClient := ConstellationClient{
			Id:               constellation.Id,
			Name:             constellation.Name,
			StartDate:        constellation.StartDate,
			EndDate:          constellation.EndDate,
			Status:           constellation.Status,
			CreationDate:     constellation.CreationDate,
			FormationDate:    constellation.FormationDate,
			ConfirmationDate: constellation.ConfirmationDate,
		}
		constellationsClient = append(constellationsClient, constellationClient)
	}
	return &constellationsClient, nil
}

func (r *Repository) GetDistinctPlanetImagesForConstellation(constellationID int) ([]PlanetImage, error) {
	var planetImages []PlanetImage
	err := r.db.Table("constellations").
		Select("DISTINCT planets.id AS planet_id, planets.name AS planet_name, planets.image_name AS image_name").
		Joins("JOIN constellations_planets ON constellations.id = constellations_planets.constellation_id").
		Joins("JOIN planets ON constellations_planets.planet_id = planets.id").
		Where("constellations.id = ?", constellationID).
		Scan(&planetImages).Error
	if err != nil {
		return nil, err
	}
	return planetImages, nil
}

/* найти созвездие по ид */
func (r *Repository) GetConstellationById(id int, userId int) (*ConstellationClientWithPlanets, error) {
	constellation := &ds.Constellation{}
	if err := r.db.Where("status != 'deleted' AND user_id= ?", userId).First(constellation, id).Error; err != nil {
		return nil, err
	}
	planets, err := r.GetDistinctPlanetImagesForConstellation(id)
	if err != nil {
		return nil, err
	}
	constellationClientWithPlanets := ConstellationClientWithPlanets{
		Id:               constellation.Id,
		Name:             constellation.Name,
		StartDate:        constellation.StartDate,
		EndDate:          constellation.EndDate,
		Status:           constellation.Status,
		CreationDate:     constellation.CreationDate,
		FormationDate:    constellation.FormationDate,
		ConfirmationDate: constellation.ConfirmationDate,
		Planets:          planets,
	}
	return &constellationClientWithPlanets, nil
}

func (r *Repository) GetConstellationByIdAdmin(id int) (*ConstellationWithPlanets, error) {
	constellation := &ds.Constellation{}
	if err := r.db.Where("status != 'deleted'").First(constellation, id).Error; err != nil {
		return nil, err
	}
	planets, err := r.GetDistinctPlanetImagesForConstellation(id)
	if err != nil {
		return nil, err
	}
	constellationWithPlanets := ConstellationWithPlanets{
		Id:               constellation.Id,
		Name:             constellation.Name,
		StartDate:        constellation.StartDate,
		EndDate:          constellation.EndDate,
		Status:           constellation.Status,
		ModeratorId:      constellation.ModeratorId,
		UserId:           constellation.UserId,
		CreationDate:     constellation.CreationDate,
		FormationDate:    constellation.FormationDate,
		ConfirmationDate: constellation.ConfirmationDate,
		Planets:          planets,
	}
	return &constellationWithPlanets, nil
}

/* возвращает последнюю заявку пользователя со статусом "created" */
func (r *Repository) GetCreatedConstellationByUser(userID int) (*ds.Constellation, error) {
	var constellation ds.Constellation
	if err := r.db.
		Where("user_id = ? AND status = 'created'", userID).
		Order("creation_date DESC").
		First(&constellation).Error; err != nil {
		return nil, err
	}
	return &constellation, nil
}

/* список всех заявок в статусе created */
func (r *Repository) GetCreatedConstellations() (*[]ds.Constellation, error) {
	var constellations []ds.Constellation

	if err := r.db.Where("status = 'created'").Order("creation_date DESC").Find(&constellations).Error; err != nil {
		return nil, err
	}
	return &constellations, nil
}

/* список всех заявок в статусе inprogress */
func (r *Repository) GetInProgressConstellations() (*[]ds.Constellation, error) {
	var constellations []ds.Constellation

	if err := r.db.Where("status = 'inprogress'").Order("creation_date DESC").Find(&constellations).Error; err != nil {
		return nil, err
	}
	return &constellations, nil
}

/* список всех заявок в статусе completed */
func (r *Repository) GetCompletedConstellations() (*[]ds.Constellation, error) {
	var constellations []ds.Constellation

	if err := r.db.Where("status = 'completed'").Order("creation_date DESC").Find(&constellations).Error; err != nil {
		return nil, err
	}
	return &constellations, nil
}

/* список всех заявок в статусе deleted */
func (r *Repository) GetDeletedConstellations() (*[]ds.Constellation, error) {
	var constellations []ds.Constellation

	if err := r.db.Where("status = 'deleted'").Order("creation_date DESC").Find(&constellations).Error; err != nil {
		return nil, err
	}
	return &constellations, nil
}

/* список всех заявок в статусе canceled */
func (r *Repository) GetCanceledConstellations() (*[]ds.Constellation, error) {
	var constellations []ds.Constellation

	if err := r.db.Where("status = 'canceled'").Order("creation_date DESC").Find(&constellations).Error; err != nil {
		return nil, err
	}
	return &constellations, nil
}

/* обновить данные заявки */
func (r *Repository) UpdateConstellationByID(id int, updatedConstellation *ds.Constellation) error {
	var constellation ds.Constellation

	constellation.Name = updatedConstellation.Name
	constellation.StartDate = updatedConstellation.StartDate
	constellation.EndDate = updatedConstellation.EndDate

	if err := r.db.Model(&ds.Constellation{}).Where("id = ?", id).Updates(constellation).Error; err != nil {
		return err
	}
	return nil
}

/* удалить созвездие */
func (r *Repository) UpdateStatusToDeleted(id int, userId int) error {
	currentTime := time.Now()
	updatedFields := ds.Constellation{
		Status:           "deleted",
		ConfirmationDate: &currentTime,
	}
	if err := r.db.Model(&ds.Constellation{}).Where("id = ? AND user_id = ?", id, userId).Updates(updatedFields).Error; err != nil {
		return err
	}
	return nil
}

/* сделать статус inprogress */
func (r *Repository) UpdateStatusToInProgress(id int, adminId uint) error {
	currentTime := time.Now()
	updatedFields := ds.Constellation{
		Status:        "inprogress",
		FormationDate: &currentTime,
	}
	if err := r.db.Model(&ds.Constellation{}).Where("id = ?", id).Updates(updatedFields).Error; err != nil {
		return err
	}
	return nil
}

/* сделать статус completed */
func (r *Repository) UpdateStatusToCompleted(id int, adminId uint) error {
	currentTime := time.Now()
	updatedFields := ds.Constellation{
		Status:           "completed",
		ConfirmationDate: &currentTime,
		ModeratorId:      &adminId,
	}
	if err := r.db.Model(&ds.Constellation{}).Where("id = ?", id).Updates(updatedFields).Error; err != nil {
		return err
	}
	return nil
}

/* сделать статус canceled */
func (r *Repository) UpdateStatusToCanceled(id int, adminId uint) error {
	currentTime := time.Now()
	updatedFields := ds.Constellation{
		Status:           "canceled",
		ConfirmationDate: &currentTime,
		ModeratorId:      &adminId,
	}
	if err := r.db.Model(&ds.Constellation{}).Where("id = ?", id).Updates(updatedFields).Error; err != nil {
		return err
	}
	return nil
}

/* создать заявку для пользователя */
func (r *Repository) CreateConstellationForUser(userId int) error {
	newConstellation := ds.Constellation{
		Name:         "Черновик",
		StartDate:    time.Now(),
		EndDate:      time.Now(),
		UserId:       uint(userId),
		Status:       "created",
		CreationDate: time.Now(),
	}
	if err := r.db.Create(&newConstellation).Error; err != nil {
		return err
	}
	return nil
}

/* удалить планету из созвездия */
func (r *Repository) DeletePlanetFromConstellation(planetID, constellationID uint) error {
	result := r.db.Exec("DELETE FROM constellations_planets WHERE planet_id = ? AND constellation_id = ?", planetID, constellationID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Удалить все связи планет из созвездия, где юзер является владельцем созвездия
func (r *Repository) DeleteAllPlanetsFromConstellation(id int, userId int) error {
	if err := r.db.Exec("DELETE FROM constellations_planets WHERE constellation_id = ? AND constellation_id IN (SELECT id FROM constellations WHERE user_id = ?)", id, userId).Error; err != nil {
		return err
	}

	return nil
}
