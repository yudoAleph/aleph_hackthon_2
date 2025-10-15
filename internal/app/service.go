package app

import (
	"context"
)

type Service struct {
	repo      UserRepository
	jwtSecret string
}

func NewService(repo UserRepository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) Get(ctx context.Context, ID string) (User, error) {
	return s.repo.Get(ctx, ID)
}
