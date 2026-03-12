package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidType   = errors.New("activity type cannot be empty")
	ErrInvalidReps   = errors.New("reps must be greater than zero")
	ErrInvalidWeight = errors.New("weight cannot be negative")
)

type Activity struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`   // например, "pull-ups"
	Weight    float64   `json:"weight"` // вес жилета или доп. отягощения
	Reps      int       `json:"reps"`   // количество повторений
	CreatedAt time.Time `json:"created_at"`
}

func (a Activity) Validate() error {
	if strings.TrimSpace(a.Type) == "" {
		return ErrInvalidType
	}
	if a.Reps <= 0 {
		return ErrInvalidReps
	}
	if a.Weight < 0 {
		return ErrInvalidWeight
	}
	return nil
}
