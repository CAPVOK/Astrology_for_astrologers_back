package model

import "time"

type Constellation struct {
	ConstellationID     uint      `gorm:"type:serial;primarykey" json:"id"`
	Name                string    `json:"name"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	CreationDate        time.Time `json:"creation_date"`
	FormationDate       time.Time `json:"formation_date"`
	CompletionDate      time.Time `json:"confirmation_date"`
	ConstellationStatus string    `json:"status"`
	UserID              uint      `json:"user_id"`
	ModeratorID         uint      `json:"moderator_id"`
}

type ConstellationRequest struct {
	DeliveryID     uint      `json:"delivery_id"`
	FlightNumber   string    `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	DeliveryStatus string    `json:"delivery_status"`
	FullName       string    `json:"full_name"`
}

type ConstellationGetResponse struct {
	DeliveryID     uint      `json:"delivery_id"`
	FlightNumber   string    `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	DeliveryStatus string    `json:"delivery_status"`
	FullName       string    `json:"full_name"`
	Baggages       []Planet  `json:"planets"`
}

type DeliveryUpdateFlightNumberRequest struct { //
	FlightNumber string `json:"flight_number"`
}

type DeliveryUpdateStatusRequest struct {
	DeliveryStatus string `json:"delivery_status"` //
}
