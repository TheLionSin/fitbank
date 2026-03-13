package app

import (
	"context"
	"fitbank/activity-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ActivityUseCase interface {
	Create(ctx context.Context, act domain.Activity) (domain.Activity, error)
	FetchAll(ctx context.Context) ([]domain.Activity, error)
	GetByID(ctx context.Context, id string) (domain.Activity, error)
	Update(ctx context.Context, act domain.Activity) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo domain.ActivityRepository
}

func NewService(repo domain.ActivityRepository) ActivityUseCase {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, act domain.Activity) (domain.Activity, error) {
	if err := act.Validate(); err != nil {
		return domain.Activity{}, err
	}

	act.ID = uuid.New().String()
	act.CreatedAt = time.Now()

	if err := s.repo.Create(ctx, act); err != nil {
		return domain.Activity{}, err
	}

	return act, nil
}

func (s *service) FetchAll(ctx context.Context) ([]domain.Activity, error) {
	return s.repo.FetchAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (domain.Activity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, act domain.Activity) error {
	// Здесь тоже важна валидация перед обновлением!
	if err := act.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, act)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
