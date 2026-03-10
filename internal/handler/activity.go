package handler

import (
	"encoding/json"
	"fitbank/internal/domain"
	"net/http"
)

type ActivityHandler struct {
}

func NewActivityHandler() *ActivityHandler {
	return &ActivityHandler{}
}

func (h *ActivityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var act domain.Activity

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	act.ID = "123"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(act)
}
