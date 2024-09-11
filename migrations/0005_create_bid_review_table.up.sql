CREATE TABLE IF NOT EXISTS bid_review (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description VARCHAR(1000) NOT NULL,
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organization(id) ON DELETE SET NULL,
    creator_id UUID NOT NULL REFERENCES employee(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
);