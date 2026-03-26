package handler

import "fitbank/activity-service/internal/domain"

// CreateActivityRequest — что мы ждем от фронтенда
type CreateActivityRequest struct {
	Type   string  `json:"type"`
	Weight float64 `json:"weight"`
	Reps   int     `json:"reps"`
}

// ToDomain — метод конвертации DTO в бизнес-сущность
func (r CreateActivityRequest) ToDomain() domain.Activity {
	return domain.Activity{
		Type:   r.Type,
		Weight: r.Weight,
		Reps:   r.Reps,
	}
}

type UpdateActivityRequest struct {
	Type   string  `json:"type"`
	Weight float64 `json:"weight"`
	Reps   int     `json:"reps"`
}

func (r UpdateActivityRequest) ToDomain(id string) domain.Activity {
	return domain.Activity{
		ID:     id,
		Type:   r.Type,
		Weight: r.Weight,
		Reps:   r.Reps,
	}
}
