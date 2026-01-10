package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"user-restriction-manager/internal/core/domain"

	"github.com/google/uuid"
)

type RestrictionRepository struct {
	db *sql.DB
}

func NewRestrictionRepository(db *sql.DB) *RestrictionRepository {
	return &RestrictionRepository{db: db}
}

func (r *RestrictionRepository) Create(ctx context.Context, res *domain.Restriction) error {
	query := `
		INSERT INTO restrictions (user_id, type, reason, start_at, end_at, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		res.UserID,
		res.Type,
		res.Reason,
		res.StartAt,
		res.EndAt,
		res.Status,
		res.CreatedBy,
	).Scan(&res.ID, &res.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create restriction: %w", err)
	}
	return nil
}

func (r *RestrictionRepository) GetActiveByUserID(ctx context.Context, userID string) ([]domain.Restriction, error) {
	query := `
		SELECT id, user_id, type, reason, start_at, end_at, status, created_by, created_at
		FROM restrictions
		WHERE user_id = $1 AND status = 'ACTIVE' AND (end_at IS NULL OR end_at > $2)`

	rows, err := r.db.QueryContext(ctx, query, userID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to query active restrictions: %w", err)
	}
	defer rows.Close()

	var restrictions []domain.Restriction
	for rows.Next() {
		var res domain.Restriction
		if err := rows.Scan(
			&res.ID,
			&res.UserID,
			&res.Type,
			&res.Reason,
			&res.StartAt,
			&res.EndAt,
			&res.Status,
			&res.CreatedBy,
			&res.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan restriction: %w", err)
		}
		restrictions = append(restrictions, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating restrictions: %w", err)
	}

	return restrictions, nil
}

func (r *RestrictionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RestrictionStatus) error {
	query := `UPDATE restrictions SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update restriction status: %w", err)
	}
	return nil
}
