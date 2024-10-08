DO $$ BEGIN
    CREATE TYPE bid_status AS ENUM ('Created', 'Published', 'Canceled', 'Approved', 'Rejected');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE bid_author_type AS ENUM ('Organization', 'User');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS bid (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    status bid_status NOT NULL,
    tender_id UUID NOT NULL,
    organization_id UUID REFERENCES organization(id) ON DELETE SET NULL,
    creator_id UUID NOT NULL REFERENCES employee(id) ON DELETE RESTRICT,
    version INT CHECK (version >= 1) DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, version)
);