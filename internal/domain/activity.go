package domain

import "time"

type Activity struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`   // например, "pull-ups"
	Weight    float64   `json:"weight"` // вес жилета или доп. отягощения
	Reps      int       `json:"reps"`   // количество повторений
	CreatedAt time.Time `json:"created_at"`
}
