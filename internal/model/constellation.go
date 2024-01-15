package model

import "time"

type Constellation struct {
	ConstellationID     uint       `gorm:"type:serial;primarykey" json:"id"`
	Name                string     `json:"name"`
	StartDate           time.Time  `json:"startDate"`
	EndDate             time.Time  `json:"endDate"`
	CreationDate        time.Time  `json:"creationDate"`
	FormationDate       *time.Time `json:"formationDate"`
	ConfirmationDate    *time.Time `json:"confirmationDate"`
	ConstellationStatus string     `json:"status"`
	UserID              uint       `json:"userId"`
	ModeratorID         *uint      `gorm:"foreignkey:constellationId" json:"moderatorId"`
}

type ConstellationRequest struct {
	ConstellationID     uint       `json:"id"`
	Name                string     `json:"name"`
	StartDate           time.Time  `json:"startDate"`
	EndDate             time.Time  `json:"endDate"`
	CreationDate        time.Time  `json:"creationDate"`
	FormationDate       *time.Time `json:"formationDate"`
	ConfirmationDate    *time.Time `json:"confirmationDate"`
	ConstellationStatus string     `json:"status"`
	FullName            string     `json:"fullName"`
}

type ConstellationGetResponse struct {
	ConstellationID     uint                    `gorm:"foreignkey:id" json:"id"`
	Name                string                  `json:"name"`
	StartDate           time.Time               `json:"startDate"`
	EndDate             time.Time               `json:"end_date"`
	CreationDate        time.Time               `json:"creationDate"`
	FormationDate       *time.Time              `json:"formationDate"`
	ConfirmationDate    *time.Time              `json:"confirmationDate"`
	ConstellationStatus string                  `json:"status"`
	FullName            string                  `json:"fullName"`
	Planets             []PlanetInConstellation `gorm:"many2many:constellationPlanets" json:"planets"`
}

type ConstellationUpdateRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type ConstellationUpdateStatusRequest struct {
	ConstellationStatus string `json:"constellationStatus"`
}
