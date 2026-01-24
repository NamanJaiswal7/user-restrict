package handler

import (
	"encoding/json"
	"net/http"

	"user-restriction-manager/internal/core/domain"
	"user-restriction-manager/internal/core/ports"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AppealHandler struct {
	service ports.AppealService
}

func NewAppealHandler(service ports.AppealService) *AppealHandler {
	return &AppealHandler{service: service}
}

type CreateAppealRequest struct {
	RestrictionID string `json:"restriction_id"`
	UserID        string `json:"user_id"`
	Reason        string `json:"reason"`
}

func (h *AppealHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req CreateAppealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resID, err := uuid.Parse(req.RestrictionID)
	if err != nil {
		http.Error(w, "invalid restriction_id", http.StatusBadRequest)
		return
	}

	appeal := &domain.Appeal{
		RestrictionID: resID,
		UserID:        req.UserID,
		Reason:        req.Reason,
	}

	if err := h.service.SubmitAppeal(r.Context(), appeal); err != nil {
		http.Error(w, "failed to submit appeal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      appeal.ID,
		"message": "appeal submitted",
	})
}

type ReviewAppealRequest struct {
	ReviewerID string `json:"reviewer_id"`
	Status     string `json:"status"` // APPROVED, REJECTED
	Notes      string `json:"notes"`
}

func (h *AppealHandler) Review(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	appealID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid appeal id", http.StatusBadRequest)
		return
	}

	var req ReviewAppealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	status := domain.AppealStatus(req.Status)
	if status != domain.AppealStatusApproved && status != domain.AppealStatusRejected {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	if err := h.service.ReviewAppeal(r.Context(), appealID, req.ReviewerID, status, req.Notes); err != nil {
		http.Error(w, "failed to review appeal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "appeal reviewed"})
}
