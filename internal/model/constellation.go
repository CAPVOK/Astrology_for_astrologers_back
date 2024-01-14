package model

import "time"

type Constellation struct {
	ConstellationID     uint       `gorm:"type:serial;primarykey" json:"id"`
	Name                string     `json:"name"`
	StartDate           time.Time  `json:"start_date"`
	EndDate             time.Time  `json:"end_date"`
	CreationDate        time.Time  `json:"creation_date"`
	FormationDate       *time.Time `json:"formation_date"`
	ConfirmationDate    *time.Time `json:"confirmation_date"`
	ConstellationStatus string     `json:"status"`
	UserID              uint       `json:"user_id"`
	ModeratorID         *uint      `gorm:"foreignkey:constellation_id" json:"moderator_id"`
}

type ConstellationRequest struct {
	ConstellationID     uint       `json:"id"`
	Name                string     `json:"name"`
	StartDate           time.Time  `json:"start_date"`
	EndDate             time.Time  `json:"end_date"`
	CreationDate        time.Time  `json:"creation_date"`
	FormationDate       *time.Time `json:"formation_date"`
	ConfirmationDate    *time.Time `json:"confirmation_date"`
	ConstellationStatus string     `json:"status"`
	FullName            string     `json:"full_name"`
}

type ConstellationGetResponse struct {
	ConstellationID     uint                    `gorm:"foreignkey:id" json:"id"`
	Name                string                  `json:"name"`
	StartDate           time.Time               `json:"start_date"`
	EndDate             time.Time               `json:"end_date"`
	CreationDate        time.Time               `json:"creation_date"`
	FormationDate       *time.Time              `json:"formation_date"`
	ConfirmationDate    *time.Time              `json:"confirmation_date"`
	ConstellationStatus string                  `json:"status"`
	FullName            string                  `json:"full_name"`
	Planets             []PlanetInConstellation `gorm:"many2many:constellation_planets" json:"planets"`
}

type ConstellationUpdateRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ConstellationUpdateStatusRequest struct {
	ConstellationStatus string `json:"constellationStatus"`
}
