package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"space/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetActivePlanets() (*[]ds.Planet, error) {
	var planets []ds.Planet

	if err := r.db.Where("is_active = ?", true).Find(&planets).Error; err != nil {
		return nil, err
	}
	return &planets, nil
}

func (r *Repository) GetActivePlanetById(id int) (*ds.Planet, error) {
	planet := &ds.Planet{}
	if err := r.db.First(planet, id).Error; err != nil {
		return nil, err
	}
	return planet, nil
}

func (r *Repository) DeactivatePlanetByID(id int) error {
	if err := r.db.Exec("UPDATE planets SET is_active=false WHERE id= ?", id).Error; err != nil {
		return err
	}
	return nil
}
