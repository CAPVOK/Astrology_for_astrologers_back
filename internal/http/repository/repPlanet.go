package repository

import (
	"errors"
	"time"

	"space/internal/model"
)

type PlanetRepository interface {
	GetPlanets(searchCode string, userID uint) (model.PlanetsGetResponse, error)
}

func (r *Repository) GetPlanets(searchName string, userID uint, page, pageSize int) (model.PlanetsGetResponse, error) {
	var constellationID uint
	if err := r.db.
		Table("constellations").
		Select("constellations.constellation_id").
		Where("user_id = ? AND constellation_status = ?", userID, model.CONSTELLATION_STATUS_DRAFT).
		Take(&constellationID).Error; err != nil {
	}

	var planets []model.Planet
	offset := (page - 1) * pageSize
	/* if err := r.db.Table("planets").
		Where("planets.planet_status = ? AND planets.name LIKE ?", model.PLANET_STATUS_ACTIVE, searchName).
		Order("planet_id").
		Scan(&planets).Error; err != nil {
		return model.PlanetsGetResponse{}, errors.New("ошибка нахождения списка планет")
	} */
	r.db.Table("planets").
		Where("planets.planet_status = ? AND planets.name LIKE ?", model.PLANET_STATUS_ACTIVE, searchName)
	r.db.Table("planets").
		Where("planets.planet_status = ? AND planets.name LIKE ?", model.PLANET_STATUS_ACTIVE, searchName)
	r.db.Table("planets").
		Where("planets.planet_status = ? AND planets.name LIKE ?", model.PLANET_STATUS_ACTIVE, searchName)
	r.db.Table("planets").
		Where("planets.planet_status = ? AND planets.name LIKE ?", model.PLANET_STATUS_ACTIVE, searchName)

	query := r.db.Table("planets").
		Where("planets.planet_status = ? AND planets.name LIKE ?", model.PLANET_STATUS_ACTIVE, searchName).
		Limit(pageSize).Offset(offset)

	if err := query.Find(&planets).Error; err != nil {
		return model.PlanetsGetResponse{}, err
	}

	planetResponse := model.PlanetsGetResponse{
		Planets:         planets,
		ConstellationID: constellationID,
	}
	return planetResponse, nil
}

func (r *Repository) GetPlanetByID(planetID, userID uint) (model.Planet, error) {
	var planet model.Planet
	if err := r.db.Table("planets").
		Where("planet_status = ? AND planet_id = ?", model.PLANET_STATUS_ACTIVE, planetID).
		First(&planet).Error; err != nil {
		return model.Planet{}, errors.New("ошибка при получении активной планеты из БД")
	}
	return planet, nil
}

func (r *Repository) CreatePlanet(userID uint, planet model.Planet) error {
	if err := r.db.Create(&planet).Error; err != nil {
		return errors.New("ошибка создания планеты")
	}
	return nil
}

func (r *Repository) DeletePlanet(planetID, userID uint) error {
	var planet model.Planet
	if err := r.db.Table("planets").
		Where("planet_id = ? AND planet_status = ?", planetID, model.PLANET_STATUS_ACTIVE).
		First(&planet).Error; err != nil {
		return errors.New("планета не найден или уже удален")
	}
	planet.PlanetStatus = model.PLANET_STATUS_DELETED
	if err := r.db.Table("planets").
		Model(&model.Planet{}).
		Where("planet_id = ?", planetID).
		Updates(planet).Error; err != nil {
		return errors.New("ошибка при обновлении статуса планеты в БД")
	}
	return nil
}

func (r *Repository) UpdatePlanet(planetID, userID uint, planet model.Planet) error {
	if err := r.db.Table("planets").
		Model(&model.Planet{}).
		Where("planet_id = ? AND planet_status = ?", planetID, model.PLANET_STATUS_ACTIVE).
		Updates(planet).Error; err != nil {
		return errors.New("ошибка при обновлении информации планеты БД")
	}

	return nil
}

func (r *Repository) AddPlanetToConstellation(planetID, userID uint) error {
	var planet model.Planet
	if err := r.db.Table("planets").
		Where("planet_id = ? AND planet_status = ?", planetID, model.PLANET_STATUS_ACTIVE).
		First(&planet).Error; err != nil {
		return errors.New("планета не найден или удален")
	}
	var constellation model.Constellation
	if err := r.db.Table("constellations").
		Where("constellation_status = ? AND user_id = ?", model.CONSTELLATION_STATUS_DRAFT, userID).
		Last(&constellation).Error; err != nil {
		constellation = model.Constellation{
			ConstellationStatus: model.CONSTELLATION_STATUS_DRAFT,
			Name:                "Созвездие",
			StartDate:           time.Now(),
			EndDate:             time.Now(),
			CreationDate:        time.Now(),
			UserID:              userID,
			ModeratorID:         nil,
		}
		if err := r.db.Table("constellations").
			Create(&constellation).Error; err != nil {
			return errors.New("ошибка создания созвездия со статусом черновик")
		}
	}
	constellationPlanet := model.ConstellationPlanet{
		PlanetID:        planetID,
		ConstellationID: constellation.ConstellationID,
	}
	if err := r.db.Table("constellation_planets").
		Create(&constellationPlanet).Error; err != nil {
		return errors.New("не удалось добавить планеты, возможно она уже добавлена")
	}
	return nil
}

func (r *Repository) RemovePlanetFromConstellation(planetID, userID uint) error {
	var constellationPlanet model.ConstellationPlanet
	if err := r.db.Joins("JOIN constellations ON constellation_planets.constellation_id = constellations.constellation_id").
		Where("constellation_planets.planet_id = ? AND constellations.user_id = ? AND constellations.constellation_status = ?", planetID, userID, model.CONSTELLATION_STATUS_DRAFT).
		First(&constellationPlanet).Error; err != nil {
		return errors.New("планета не принадлежит пользователю или находится не в статусе черновик")
	}
	if err := r.db.Table("constellation_planets").
		Delete(&constellationPlanet).Error; err != nil {
		return errors.New("ошибка удаления связи между созвездием и планетой")
	}
	return nil
}

func (r *Repository) AddPlanetImage(planetID, userID uint, imageURL string) error {
	err := r.db.Table("planets").Where("planet_id = ?", planetID).Update("image_name", imageURL).Error
	if err != nil {
		return errors.New("ошибка обновления url изображения в БД")
	}

	return nil
}
