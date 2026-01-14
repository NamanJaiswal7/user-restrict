CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE restriction_type AS ENUM ('WARNING', 'TEMP_BAN', 'PERM_BAN');
CREATE TYPE restriction_status AS ENUM ('ACTIVE', 'EXPIRED', 'REVOKED');
CREATE TYPE appeal_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TABLE restrictions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    type restriction_type NOT NULL,
    reason TEXT NOT NULL,
    start_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    end_at TIMESTAMP WITH TIME ZONE, -- NULL for perm ban
    status restriction_status NOT NULL DEFAULT 'ACTIVE',
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_restrictions_user_id ON restrictions(user_id);
CREATE INDEX idx_restrictions_status ON restrictions(status);

CREATE TABLE appeals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    restriction_id UUID NOT NULL REFERENCES restrictions(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    reason TEXT NOT NULL,
    status appeal_status NOT NULL DEFAULT 'PENDING',
    reviewer_id VARCHAR(255),
    review_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_appeals_user_id ON appeals(user_id);
CREATE INDEX idx_appeals_restriction_id ON appeals(restriction_id);
CREATE INDEX idx_appeals_status ON appeals(status);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    action VARCHAR(255) NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    actor_id VARCHAR(255) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
