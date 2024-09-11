DO $$ BEGIN
    CREATE TYPE bid_decision_type AS ENUM ('Approved', 'Rejected');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;


CREATE TABLE IF NOT EXISTS bid_decision (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bid_id UUID NOT NULL,
    type bid_decision_type NOT NULL, 
    organization_id UUID REFERENCES organization(id) ON DELETE SET NULL,
    creator_id UUID NOT NULL REFERENCES employee(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);