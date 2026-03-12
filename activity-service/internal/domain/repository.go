package domain

import "context"

type ActivityRepository interface {
	Create(ctx context.Context, act Activity) error
	FetchAll(ctx context.Context) ([]Activity, error)
	GetByID(ctx context.Context, id string) (Activity, error)
	Update(ctx context.Context, act Activity) error
	Delete(ctx context.Context, id string) error
}
