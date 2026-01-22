package ports

import (
	"context"

	"user-restriction-manager/internal/core/domain"

	"github.com/google/uuid"
)

type RestrictionRepository interface {
	Create(ctx context.Context, res *domain.Restriction) error
	GetActiveByUserID(ctx context.Context, userID string) ([]domain.Restriction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RestrictionStatus) error
}

type AppealRepository interface {
	Create(ctx context.Context, appeal *domain.Appeal) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Appeal, error)
	Update(ctx context.Context, appeal *domain.Appeal) error
}

type CacheRepository interface {
	CacheActiveRestrictions(ctx context.Context, userID string, restrictions []domain.Restriction) error
	GetActiveRestrictions(ctx context.Context, userID string) ([]domain.Restriction, error)
	Invalidate(ctx context.Context, userID string) error
}

type RestrictionService interface {
	ApplyRestriction(ctx context.Context, res *domain.Restriction) error
	GetActiveRestrictions(ctx context.Context, userID string) ([]domain.Restriction, error)
	RevokeRestriction(ctx context.Context, id uuid.UUID, reason string) error
}

type AppealService interface {
	SubmitAppeal(ctx context.Context, appeal *domain.Appeal) error
	ReviewAppeal(ctx context.Context, appealID uuid.UUID, reviewerID string, status domain.AppealStatus, notes string) error
}
