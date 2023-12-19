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

//type User struct {
//	Id       uint `gorm:"primaryKey"`
//	Login    string
//	Password string
//	Admin    bool
//	User     []User
//}

type User struct {
	ID       uint   `gorm:"primarykey;autoIncrement"`
	Email    string `gorm:"type:varchar(30); unique"`
	Password []byte `gorm:"type:bytea" json:"password,omitempty"`
	ImageRef string `gorm:"type:varchar(90)" json:"imageRef,omitempty"Z`
	Role     Role   `gorm:"type:int;" json:"role,omitempty"`
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

type ConstellationsRequest struct {
	Id               uint `gorm:"primaryKey"`
	Name             string
	StartDate        time.Time
	EndDate          time.Time
	Status           string
	CreationDate     time.Time
	FormationDate    *time.Time
	ConfirmationDate *time.Time
	UserEmail		 string
}

type ConstellationsPlanets struct {
	PlanetID        uint `gorm:"primaryKey;index"`
	ConstellationID uint `gorm:"primaryKey;index"`
}
