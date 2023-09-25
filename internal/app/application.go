package app

import (
	"space/internal/app/dsn"
	"space/internal/app/repository"

	"github.com/joho/godotenv"
)

type Application struct {
	repo *repository.Repository
}

func New() (Application, error) {
	_ = godotenv.Load()
	repository, err := repository.New(dsn.FromEnv())
	if err != nil {
		return Application{}, err
	}

	return Application{repo: repository}, nil
}
