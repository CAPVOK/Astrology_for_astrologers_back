package app

import (
	"context"

	"space/internal/config"
	"space/internal/dsn"
	"space/internal/http/delivery"
	"space/internal/http/repository"
	"space/internal/http/usecase"
)

// Application представляет основное приложение.
type Application struct {
	Config     *config.Config
	Repository *repository.Repository
	UseCase    *usecase.UseCase
	Handler    *delivery.Handler
}

// New создает новый объект Application и настраивает его.
func New(ctx context.Context) (*Application, error) {
	// Инициализируйте конфигурацию
	cfg, err := config.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Инициализируйте подключение к базе данных (DB)
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}
	uc := usecase.NewUseCase(repo)
	h := delivery.NewHandler(uc)
	// Инициализируйте и настройте объект Application
	app := &Application{
		Config:     cfg,
		Repository: repo,
		UseCase:    uc,
		Handler:    h,
	}

	return app, nil
}
