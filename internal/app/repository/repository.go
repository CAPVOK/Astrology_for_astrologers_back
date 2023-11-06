package repository

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"space/internal/app/minioclient"
)

type Repository struct {
	db          *gorm.DB
	minioClient *minioclient.MinioClient
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	minioClient, err := minioclient.NewMinioClient()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:          db,
		minioClient: minioClient,
	}, nil
}

type PlanetClient struct {
	Id         uint `gorm:"primarykey:autoIncrement"`
	Name       string
	Discovered *string
	Mass       *string
	Distance   *string
	Info       *string
	Color1     *string
	Color2     *string
	ImageName  *string
}

type PlanetImage struct {
	PlanetID   uint `gorm:"primarykey:autoIncrement"`
	PlanetName string
	ImageName  string
}

type ConstellationWithPlanets struct {
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
	Planets          []PlanetImage
}

type ConstellationClient struct {
	Id               uint
	Name             string
	StartDate        time.Time
	EndDate          time.Time
	Status           string
	CreationDate     time.Time
	FormationDate    *time.Time
	ConfirmationDate *time.Time
}

type ConstellationClientWithPlanets struct {
	Id               uint
	Name             string
	StartDate        time.Time
	EndDate          time.Time
	Status           string
	CreationDate     time.Time
	FormationDate    *time.Time
	ConfirmationDate *time.Time
	Planets          []PlanetImage
}
