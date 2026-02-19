package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SEIEKSHION/Exchanger/internal/models"
	"github.com/SEIEKSHION/Exchanger/internal/repository"
)

type MeasurementHandler struct {
	repo *repository.MeasurementRepository
}

func NewMeasurementHandler(repo *repository.MeasurementRepository) *MeasurementHandler {
	return &MeasurementHandler{repo: repo}
}

func (h *MeasurementHandler) Create(w http.ResponseWriter, r *http.Request) {
	var m models.Measurement
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Устанавливаем текущую дату, если не указана
	if m.Date.IsZero() {
		m.Date = time.Now()
	}

	if err := h.repo.Create(&m); err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
