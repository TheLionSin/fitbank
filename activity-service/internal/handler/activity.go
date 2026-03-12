package handler

import (
	"encoding/json"
	"fitbank/activity-service/internal/domain"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ActivityHandler struct {
	repo domain.ActivityRepository
}

func NewActivityHandler(repo domain.ActivityRepository) *ActivityHandler {
	return &ActivityHandler{
		repo: repo,
	}
}

func (h *ActivityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var act domain.Activity

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := act.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	act.ID = uuid.New().String()
	act.CreatedAt = time.Now()

	if err := h.repo.Create(r.Context(), act); err != nil {
		log.Printf("DB error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(act)
}

func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) {
	activities, err := h.repo.FetchAll(r.Context())
	if err != nil {
		log.Printf("failed to fetch activities: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activities)
}

func (h *ActivityHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	act, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get activity", "id", id, "error", err)
		http.Error(w, "activity not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(act)
}

func (h *ActivityHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var act domain.Activity
	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	act.ID = id

	if err := act.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := h.repo.Update(r.Context(), act); err != nil {
		if err.Error() == "activity not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		slog.Error("failed to update activity", "id", id, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(act)
}

func (h *ActivityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
