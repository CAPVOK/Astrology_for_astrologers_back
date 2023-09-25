package ds

import "time"

type Planet struct {
	Id             uint `gorm:"primarykey"`
	Name           string
	Discovered     string
	Mass           string
	Distance       string
	Info           string
	Color1         string
	Color2         string
	ImageName      string
	IsActive       bool
	Constellations []Constellation `gorm:"many2many:constellation_planet"`
}

type User struct {
	Id       uint `gorm:"primaryKey"`
	Login    string
	Password string
	Admin    bool
	User     []User
}

type Constellation struct {
	Id               uint `gorm:"primaryKey"`
	Name             string
	StartDate        time.Time
	EndDate          time.Time
	ModeratorId      uint
	UserId           uint
	Status           string
	CreationDate     time.Time
	FormationDate    time.Time
	ConfirmationDate time.Time
	Planets          []Planet `gorm:"many2many:constellation_planet"`
}
