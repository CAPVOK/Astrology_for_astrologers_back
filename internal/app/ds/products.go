package ds

import "time"

type Planet struct {
	Id             uint `gorm:"primarykey:autoIncrement"`
	Name           string
	Discovered     *string         `gorm:"default:'Неизвестно'"`
	Mass           *string         `gorm:"default:'Неизвестно'"`
	Distance       *string         `gorm:"default:'Неизвестно'"`
	Info           *string         `gorm:"default:'Неизвестно'"`
	Color1         *string         `gorm:"default:'#ababab'"`
	Color2         *string         `gorm:"default:'#8a8a8a'"`
	ImageName      *string         `gorm:"default:'unknown.png'"`
	IsActive       *bool           `gorm:"default:true"`
	Constellations []Constellation `gorm:"many2many:constellations_planets"`
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
	ModeratorId      *uint
	UserId           uint
	Status           string
	CreationDate     time.Time
	FormationDate    *time.Time
	ConfirmationDate *time.Time
	Planets          []Planet `gorm:"many2many:constellations_planets"`
}

type ConstellationsPlanets struct {
	PlanetID        uint `gorm:"primaryKey;index"`
	ConstellationID uint `gorm:"primaryKey;index"`
}
