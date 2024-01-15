package model

type ConstellationPlanet struct {
	ConstellationID uint `gorm:"type:serial;primaryKey;index" json:"constellationId"`
	PlanetID        uint `gorm:"type:serial;primaryKey;index" json:"planetId"`
}
