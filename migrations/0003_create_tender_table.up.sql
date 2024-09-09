CREATE TYPE IF NOT EXISTS tender_service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');
CREATE TYPE IF NOT EXISTS tender_status AS ENUM ('Created', 'Published', 'Closed');
CREATE TABLE IF NOT EXISTS tender (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    service_type tender_service_type NOT NULL,
    status tender_status NOT NULL,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE RESTRICT,
    creator_id UUID NOT NULL REFERENCES employee(id) ON DELETE RESTRICT,
    version INT CHECK (version >= 1) DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, version)
);