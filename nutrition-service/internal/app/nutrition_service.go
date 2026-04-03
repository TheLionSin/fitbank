package app

import (
	"context"
	"fitbank/nutrition-service/internal/client"
	"fitbank/nutrition-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type NutritionUseCase interface {
	AddMeal(ctx context.Context, meal domain.Meal) (domain.Meal, error)
	GetDailyReport(ctx context.Context) (domain.DailyReport, error)
}

type service struct {
	repo           domain.MealRepository
	activityClient client.ActivityClientInterface
}

func NewService(repo domain.MealRepository, actClient client.ActivityClientInterface) NutritionUseCase {
	return &service{
		repo:           repo,
		activityClient: actClient,
	}
}

func (s *service) AddMeal(ctx context.Context, meal domain.Meal) (domain.Meal, error) {
	if err := meal.Validate(); err != nil {
		return domain.Meal{}, err
	}

	meal.ID = uuid.New().String()
	meal.CreatedAt = time.Now()

	if err := s.repo.Create(ctx, meal); err != nil {
		return domain.Meal{}, err
	}

	return meal, nil
}

// GetDailyReport — создание отчета за день
func (s *service) GetDailyReport(ctx context.Context) (domain.DailyReport, error) {
	// 1. Получаем сумму съеденного из нашей БД
	totals, err := s.repo.GetDailyTotal(ctx)
	if err != nil {
		return domain.DailyReport{}, err
	}

	// 2. Звоним в activity-service за сожженными калориями
	burned, err := s.activityClient.GetTotalCaloriesBurned(ctx)
	if err != nil {
		return domain.DailyReport{}, err
	}

	// 3. Формируем итоговый отчет
	return domain.DailyReport{
		TotalCalories:  totals.Calories,
		TotalProteins:  totals.Proteins,
		TotalFats:      totals.Fats,
		TotalCarbs:     totals.Carbohydrates,
		BurnedCalories: burned,
		NetCalories:    totals.Calories - burned, // Чистый итог (съедено - потрачено)
	}, nil
}
