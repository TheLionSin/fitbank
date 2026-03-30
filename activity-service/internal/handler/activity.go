package handler

import (
	"encoding/json"
	"fitbank/activity-service/internal/app"
	"fitbank/activity-service/internal/metrics"
	"log"
	"log/slog"
	"net/http"
)

type ActivityHandler struct {
	service app.ActivityUseCase
}

func NewActivityHandler(service app.ActivityUseCase) *ActivityHandler {
	return &ActivityHandler{
		service: service,
	}
}

func (h *ActivityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Вызываем сервис. Он сам провалидирует, создаст UUID и сохранит.
	// Маппим DTO в доменную модель и отдаем сервису
	act, err := h.service.Create(r.Context(), req.ToDomain())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	metrics.ActivitiesCreatedTotal.Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(act)
}

func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) {
	activities, err := h.service.FetchAll(r.Context())
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

	act, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get activity", "id", id, "error", err)
		http.Error(w, "activity not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(act)
}

func (h *ActivityHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req UpdateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), req.ToDomain(id)); err != nil {
		if err.Error() == "activity not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ActivityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
