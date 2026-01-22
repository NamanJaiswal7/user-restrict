package service

import (
	"context"
	"fmt"
	"time"

	"user-restriction-manager/internal/core/domain"
	"user-restriction-manager/internal/core/ports"

	"github.com/google/uuid"
)

type AppealService struct {
	appealRepo      ports.AppealRepository
	restrictionRepo ports.RestrictionRepository
	// notificationService ports.NotificationService // Future
}

func NewAppealService(appealRepo ports.AppealRepository, restrictionRepo ports.RestrictionRepository) *AppealService {
	return &AppealService{
		appealRepo:      appealRepo,
		restrictionRepo: restrictionRepo,
	}
}

func (s *AppealService) SubmitAppeal(ctx context.Context, appeal *domain.Appeal) error {
	// Validate if restriction exists and is active?
	// Assuming validation happens or DB constraints catch it.
	
	appeal.Status = domain.AppealStatusPending
	appeal.CreatedAt = time.Now()
	appeal.UpdatedAt = time.Now()

	if err := s.appealRepo.Create(ctx, appeal); err != nil {
		return fmt.Errorf("service: failed to submit appeal: %w", err)
	}
	return nil
}

func (s *AppealService) ReviewAppeal(ctx context.Context, appealID uuid.UUID, reviewerID string, status domain.AppealStatus, notes string) error {
	appeal, err := s.appealRepo.GetByID(ctx, appealID)
	if err != nil {
		return fmt.Errorf("service: failed to fetch appeal: %w", err)
	}
	if appeal == nil {
		return fmt.Errorf("appeal not found")
	}

	appeal.Status = status
	appeal.ReviewerID = &reviewerID
	appeal.ReviewNotes = &notes
	appeal.UpdatedAt = time.Now()

	if err := s.appealRepo.Update(ctx, appeal); err != nil {
		return fmt.Errorf("service: failed to update appeal: %w", err)
	}

	// If approved, revoke restriction
	if status == domain.AppealStatusApproved {
		if err := s.restrictionRepo.UpdateStatus(ctx, appeal.RestrictionID, domain.RestrictionStatusRevoked); err != nil {
			return fmt.Errorf("service: failed to revoke restriction after appeal approval: %w", err)
		}
	}

	return nil
}
