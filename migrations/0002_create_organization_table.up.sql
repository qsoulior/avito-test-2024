CREATE TYPE IF NOT EXISTS organization_type AS ENUM ('IE', 'LLC', 'JSC');
CREATE TABLE IF NOT EXISTS organization (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE organization_responsible (
    id SERIAL PRIMARY KEY,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);
