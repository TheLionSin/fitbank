package repository

import (
	"context"
	"fitbank/activity-service/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, act domain.Activity) error {
	query := `INSERT INTO activities (id,type,weight,reps,created_at)
			VALUES ($1,$2,$3,$4,$5)`

	_, err := r.db.Exec(ctx, query, act.ID, act.Type, act.Weight, act.Reps, act.CreatedAt)
	return err
}

func (r *PostgresRepository) FetchAll(ctx context.Context) ([]domain.Activity, error) {
	query := `SELECT id,type,weight,reps,created_at FROM activities ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []domain.Activity

	for rows.Next() {
		var act domain.Activity
		err := rows.Scan(&act.ID, &act.Type, &act.Weight, &act.Reps, &act.CreatedAt)
		if err != nil {
			return nil, err
		}
		activities = append(activities, act)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if activities == nil {
		activities = []domain.Activity{}
	}

	return activities, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (domain.Activity, error) {
	var act domain.Activity
	query := `SELECT id,type,weight,reps,created_at FROM activities WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(&act.ID, &act.Type, &act.Weight, &act.Reps, &act.CreatedAt)
	return act, err
}

func (r *PostgresRepository) Update(ctx context.Context, act domain.Activity) error {
	query := `UPDATE activities SET type = $1, weight = $2, reps = $3 WHERE id = $4`

	res, err := r.db.Exec(ctx, query, act.Type, act.Weight, act.Reps, act.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("activity not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM activities WHERE id = $1`

	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("activity not found")
	}
	return nil
}
