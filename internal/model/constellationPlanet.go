package model

type ConstellationPlanet struct {
	ConstellationID uint `gorm:"type:serial;primaryKey;index" json:"constellation_id"`
	PlanetID        uint `gorm:"type:serial;primaryKey;index" json:"planet_id"`
}
