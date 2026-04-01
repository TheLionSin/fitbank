package domain

import (
	"context"
	"errors"
	"time"
)

type Meal struct {
	ID            string
	Name          string
	Calories      float64
	Proteins      float64
	Fats          float64
	Carbohydrates float64
	WeightGrams   int
	CreatedAt     time.Time
}

func (m *Meal) Validate() error {
	if m.Name == "" {
		return errors.New("meal name is required")
	}
	if m.WeightGrams <= 0 {
		return errors.New("weight must be greater than 0")
	}
	if m.Calories < 0 || m.Proteins < 0 {
		return errors.New("nutrients cannot be negative")
	}
	return nil
}

// Интерфейс для работы с БД (Репозиторий)
type MealRepository interface {
	Create(ctx context.Context, meal Meal) error
	GetDailyTotal(ctx context.Context) (Meal, error) // Возвращает сумму БЖУ за сегодня
}
