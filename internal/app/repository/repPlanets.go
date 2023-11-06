package repository

import "space/internal/app/ds"

/* планеты */

/* список активных планет */
func (r *Repository) GetActivePlanets() (*[]PlanetClient, error) {
	var planets []ds.Planet
	if err := r.db.Where("is_active = ?", true).Order("id").Find(&planets).Error; err != nil {
		return nil, err
	}
	/* return &planets, nil */

	var newPlanets []PlanetClient
	for _, planet := range planets {
		newPlanet := PlanetClient{
			Id:         planet.Id,
			Name:       planet.Name,
			Discovered: planet.Discovered,
			Mass:       planet.Mass,
			Distance:   planet.Distance,
			Info:       planet.Info,
			Color1:     planet.Color1,
			Color2:     planet.Color2,
			ImageName:  planet.ImageName,
		}
		newPlanets = append(newPlanets, newPlanet)
	}
	return &newPlanets, nil
}

/* список планет для админа */
func (r *Repository) GetPlanets() (*[]ds.Planet, error) {
	var planets []ds.Planet

	if err := r.db.Where("is_active = ?", true).Order("id").Find(&planets).Error; err != nil {
		return nil, err
	}
	return &planets, nil
}

/* активная планета по ид */
func (r *Repository) GetActivePlanetById(id int) (*PlanetClient, error) {
	planet := &ds.Planet{}
	if err := r.db.Where("is_active = ?", true).First(planet, id).Error; err != nil {
		return nil, err
	}
	/* return planet, nil */
	newPlanet := PlanetClient{
		Id:         planet.Id,
		Name:       planet.Name,
		Discovered: planet.Discovered,
		Mass:       planet.Mass,
		Distance:   planet.Distance,
		Info:       planet.Info,
		Color1:     planet.Color1,
		Color2:     planet.Color2,
		ImageName:  planet.ImageName,
	}
	return &newPlanet, nil
}

/* планета по ид */
func (r *Repository) GetPlanetById(id int) (*ds.Planet, error) {
	planet := &ds.Planet{}
	if err := r.db.First(planet, id).Error; err != nil {
		return nil, err
	}
	return planet, nil
}

/* найти планеты по имени */
func (r *Repository) SearchPlanetsByName(name string) (*[]PlanetClient, error) {
	var foundPlanets []ds.Planet
	if err := r.db.Where("is_active = ? AND LOWER(name) LIKE LOWER(?)", true, name+"%").Order("id").Find(&foundPlanets).Error; err != nil {
		return nil, err
	}
	/* return foundPlanets, nil */
	var newPlanets []PlanetClient
	for _, planet := range foundPlanets {
		newPlanet := PlanetClient{
			Id:         planet.Id,
			Name:       planet.Name,
			Discovered: planet.Discovered,
			Mass:       planet.Mass,
			Distance:   planet.Distance,
			Info:       planet.Info,
			Color1:     planet.Color1,
			Color2:     planet.Color2,
			ImageName:  planet.ImageName,
		}
		newPlanets = append(newPlanets, newPlanet)
	}
	return &newPlanets, nil
}

/* найти планеты по имени админ*/
func (r *Repository) SearchPlanetsByNameAdmin(name string) ([]ds.Planet, error) {
	var foundPlanets []ds.Planet
	if err := r.db.Where("is_active = ? AND LOWER(name) LIKE LOWER(?)", true, name+"%").Order("id").Find(&foundPlanets).Error; err != nil {
		return nil, err
	}
	return foundPlanets, nil
}

/* удалить планету */
func (r *Repository) DeactivatePlanetByID(id int) error {
	err := r.minioClient.RemoveServiceImage(id)
	if err != nil {
		return err
	}
	if err := r.db.Exec("UPDATE planets SET is_active=false WHERE id= ?", id).Error; err != nil {
		return err
	}
	err = r.minioClient.RemoveServiceImage(id)
	if err != nil {
		// Обработка ошибки удаления изображения из MinIO, если необходимо
		return err
	}
	return nil
}

/* обновить данные планеты */
func (r *Repository) UpdatePlanetByID(id int, updatedPlanet *ds.Planet) error {
	if err := r.db.Model(&ds.Planet{}).Where("id = ?", id).Updates(updatedPlanet).Error; err != nil {
		return err
	}
	return nil
}

/* создать планету */
func (r *Repository) CreatePlanet(planet *ds.Planet) error {
	if err := r.db.Create(planet).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddPlanetImage(planetId int, imageBytes []byte, contentType string) error {
	// Удаление существующего изображения (если есть)
	err := r.minioClient.RemoveServiceImage(planetId)
	if err != nil {
		return err
	}
	// Загрузка нового изображения в MinIO
	imageURL, err := r.minioClient.UploadServiceImage(planetId, imageBytes, contentType)
	if err != nil {
		return err
	}
	// Обновление информации об изображении в БД (например, ссылки на MinIO)
	err = r.db.Model(&ds.Planet{}).Where("id = ?", planetId).Update("image_name", imageURL).Error
	if err != nil {
		// Обработка ошибки обновления URL изображения в БД, если необходимо
		return err
	}
	return nil
}
