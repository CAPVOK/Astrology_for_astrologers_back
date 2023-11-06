package repository

import "space/internal/app/ds"

/* m-m */

/* добавляет планету к заявке (созвездию) по их идентификаторам */
func (r *Repository) AddPlanetToConstellation(planetID, constellationID uint) error {
	constPlanets := ds.ConstellationsPlanets{
		PlanetID:        planetID,
		ConstellationID: constellationID,
	}
	if err := r.db.Create(&constPlanets).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllConstPlanets() (*[]ds.ConstellationsPlanets, error) {
	var constplanets []ds.ConstellationsPlanets

	if err := r.db.Find(&constplanets).Error; err != nil {
		return nil, err
	}
	return &constplanets, nil
}

func (r *Repository) GetActiveConstPlanets() (*[]ds.ConstellationsPlanets, error) {
	var constplanets []ds.ConstellationsPlanets
	if err := r.db.Joins("JOIN constellations ON constellations_planets.constellation_id = constellations.id").
		Where("constellations.status != 'deleted'").
		Find(&constplanets).Error; err != nil {
		return nil, err
	}
	return &constplanets, nil
}

func (r *Repository) RemovePlanetFromConstellation(planetID, constellationID int) error {
	if err := r.db.Exec("DELETE FROM constellations_planets WHERE planet_id = ? AND constellation_id = ?", planetID, constellationID).Error; err != nil {
		return err
	}
	return nil
}
