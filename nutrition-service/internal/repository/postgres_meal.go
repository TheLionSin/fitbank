package repository

import (
	"context"
	"fitbank/nutrition-service/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMealRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMealRepository(db *pgxpool.Pool) *PostgresMealRepository {
	return &PostgresMealRepository{
		db: db,
	}
}

func (r *PostgresMealRepository) Create(ctx context.Context, m domain.Meal) error {
	query := `INSERT INTO meals (id,name,calories,proteins,fats,carbohydrates,weight_grams,created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := r.db.Exec(ctx, query,
		m.ID,
		m.Name,
		m.Calories,
		m.Proteins,
		m.Fats,
		m.Carbohydrates,
		m.WeightGrams,
		m.CreatedAt,
	)
	return err
}

// GetDailyTotal — это "мозг" нашего мониторинга питания.
// Он суммирует все показатели за текущие сутки (от 00:00:00 до 23:59:59).
func (r *PostgresMealRepository) GetDailyTotal(ctx context.Context) (domain.Meal, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	query := `SELECT COALESCE(SUM(calories),0),
					COALESCE(SUM(proteins),0),
					COALESCE(SUM(fats),0),
					COALESCE(SUM(carbohydrates),0),
			FROM meals
			WHERE created_at >= $1`

	var total domain.Meal
	err := r.db.QueryRow(ctx, query, startOfDay).Scan(
		&total.Calories,
		&total.Proteins,
		&total.Fats,
		&total.Carbohydrates)

	if err != nil {
		return domain.Meal{}, err
	}

	return total, nil
}
