package domain

import (
	"time"

	"github.com/google/uuid"
)

type RestrictionType string

const (
	RestrictionTypeWarning RestrictionType = "WARNING"
	RestrictionTypeTempBan RestrictionType = "TEMP_BAN"
	RestrictionTypePermBan RestrictionType = "PERM_BAN"
)

type RestrictionStatus string

const (
	RestrictionStatusActive  RestrictionStatus = "ACTIVE"
	RestrictionStatusExpired RestrictionStatus = "EXPIRED"
	RestrictionStatusRevoked RestrictionStatus = "REVOKED"
)

type AppealStatus string

const (
	AppealStatusPending  AppealStatus = "PENDING"
	AppealStatusApproved AppealStatus = "APPROVED"
	AppealStatusRejected AppealStatus = "REJECTED"
)

// Restriction represents a user restriction rule
type Restriction struct {
	ID        uuid.UUID         `json:"id" db:"id"`
	UserID    string            `json:"user_id" db:"user_id"`
	Type      RestrictionType   `json:"type" db:"type"`
	Reason    string            `json:"reason" db:"reason"`
	StartAt   time.Time         `json:"start_at" db:"start_at"`
	EndAt     *time.Time        `json:"end_at,omitempty" db:"end_at"`
	Status    RestrictionStatus `json:"status" db:"status"`
	CreatedBy string            `json:"created_by" db:"created_by"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
}

// Appeal represents a user's appeal against a restriction
type Appeal struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	RestrictionID uuid.UUID    `json:"restriction_id" db:"restriction_id"`
	UserID        string       `json:"user_id" db:"user_id"`
	Reason        string       `json:"reason" db:"reason"`
	Status        AppealStatus `json:"status" db:"status"`
	ReviewerID    *string      `json:"reviewer_id,omitempty" db:"reviewer_id"`
	ReviewNotes   *string      `json:"review_notes,omitempty" db:"review_notes"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

// AuditLog tracks important system actions
type AuditLog struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Action     string    `json:"action" db:"action"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	ActorID    string    `json:"actor_id" db:"actor_id"`
	Metadata   map[string]interface{} `json:"metadata" db:"metadata"` // Requires custom JSONB handling likely
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
