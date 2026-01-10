package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"user-restriction-manager/internal/core/domain"

	"github.com/google/uuid"
)

type AppealRepository struct {
	db *sql.DB
}

func NewAppealRepository(db *sql.DB) *AppealRepository {
	return &AppealRepository{db: db}
}

func (r *AppealRepository) Create(ctx context.Context, appeal *domain.Appeal) error {
	query := `
		INSERT INTO appeals (restriction_id, user_id, reason, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		appeal.RestrictionID,
		appeal.UserID,
		appeal.Reason,
		appeal.Status,
	).Scan(&appeal.ID, &appeal.CreatedAt, &appeal.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create appeal: %w", err)
	}
	return nil
}

func (r *AppealRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appeal, error) {
	query := `
		SELECT id, restriction_id, user_id, reason, status, reviewer_id, review_notes, created_at, updated_at
		FROM appeals
		WHERE id = $1`

	var appeal domain.Appeal
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&appeal.ID,
		&appeal.RestrictionID,
		&appeal.UserID,
		&appeal.Reason,
		&appeal.Status,
		&appeal.ReviewerID,
		&appeal.ReviewNotes,
		&appeal.CreatedAt,
		&appeal.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get appeal: %w", err)
	}
	return &appeal, nil
}

func (r *AppealRepository) Update(ctx context.Context, appeal *domain.Appeal) error {
	query := `
		UPDATE appeals
		SET status = $1, reviewer_id = $2, review_notes = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		appeal.Status,
		appeal.ReviewerID,
		appeal.ReviewNotes,
		appeal.ID,
	).Scan(&appeal.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update appeal: %w", err)
	}
	return nil
}
