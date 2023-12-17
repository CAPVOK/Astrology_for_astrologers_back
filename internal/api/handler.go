package api

import (
	"space/internal/app/repository"
)

type Handler struct {
	Repo      *repository.Repository
	RedisRepo repository.Redis
}

func NewHandler(repo *repository.Repository, redisRepo repository.Redis) *Handler {
	return &Handler{
		Repo:      repo,
		RedisRepo: redisRepo,
	}
}
