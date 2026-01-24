package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"user-restriction-manager/internal/core/domain"
	"user-restriction-manager/internal/core/ports"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RestrictionHandler struct {
	service ports.RestrictionService
}

func NewRestrictionHandler(service ports.RestrictionService) *RestrictionHandler {
	return &RestrictionHandler{service: service}
}

type CreateRestrictionRequest struct {
	UserID    string           `json:"user_id"`
	Type      string           `json:"type"`
	Reason    string           `json:"reason"`
	Duration  string           `json:"duration"` // ISO 8601 duration or simple string like "24h"
	CreatedBy string           `json:"created_by"`
}

func (h *RestrictionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRestrictionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate inputs (simple validation)
	if req.UserID == "" || req.Type == "" {
		http.Error(w, "user_id and type are required", http.StatusBadRequest)
		return
	}

	res := &domain.Restriction{
		UserID:    req.UserID,
		Type:      domain.RestrictionType(req.Type),
		Reason:    req.Reason,
		CreatedBy: req.CreatedBy,
	}

	if req.Type != string(domain.RestrictionTypePermBan) {
		duration, err := time.ParseDuration(req.Duration)
		if err != nil {
			// Fallback: Default 24h if invalid
			duration = 24 * time.Hour
		}
		endAt := time.Now().Add(duration)
		res.EndAt = &endAt
	}

	res.StartAt = time.Now()

	if err := h.service.ApplyRestriction(r.Context(), res); err != nil {
		http.Error(w, "failed to apply restriction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      res.ID,
		"message": "restriction applied",
	})
}

func (h *RestrictionHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	restrictions, err := h.service.GetActiveRestrictions(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch restrictions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restrictions)
}

func (h *RestrictionHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid restriction id", http.StatusBadRequest)
		return
	}

	// In real app, read reason from body
	if err := h.service.RevokeRestriction(r.Context(), id, "API Revoke"); err != nil {
		http.Error(w, "failed to revoke restriction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "restriction revoked"})
}
