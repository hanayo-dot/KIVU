-- 0001_init_postgis_schema.sql
-- Enables PostGIS, UUID extensions, and creates core KIVU entities with JWT Auth

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. Farmers Table (With Bcrypt Password Hash)
CREATE TABLE IF NOT EXISTS farmers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    location_name VARCHAR(255) NOT NULL DEFAULT 'Lake Victoria Basin',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Refresh Tokens Table (Hashed Refresh Tokens for Token Rotation)
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    farmer_id UUID NOT NULL REFERENCES farmers(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash);

-- 3. Farms Table
CREATE TABLE IF NOT EXISTS farms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    farmer_id UUID NOT NULL REFERENCES farmers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. Cages Table (PostGIS Geography Point)
CREATE TABLE IF NOT EXISTS cages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    farm_id UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    latitude NUMERIC(10, 6) NOT NULL,
    longitude NUMERIC(10, 6) NOT NULL,
    location GEOGRAPHY(Point, 4326),
    installed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active'
);

CREATE INDEX IF NOT EXISTS idx_cages_location ON cages USING GIST(location);

-- 5. Readings Table (Water Quality Measurements)
CREATE TABLE IF NOT EXISTS readings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cage_id UUID NOT NULL REFERENCES cages(id) ON DELETE CASCADE,
    dissolved_oxygen NUMERIC(5, 2) NOT NULL, -- mg/L
    temperature NUMERIC(5, 2) NOT NULL,       -- °C
    ph NUMERIC(4, 2) NOT NULL,
    turbidity NUMERIC(6, 2) NOT NULL,        -- NTU
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    source VARCHAR(20) DEFAULT 'manual'      -- 'sensor' | 'manual'
);

CREATE INDEX IF NOT EXISTS idx_readings_cage_recorded ON readings(cage_id, recorded_at DESC);

-- 6. Observations Table (Qualitative Farmer Reports)
CREATE TABLE IF NOT EXISTS observations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cage_id UUID NOT NULL REFERENCES cages(id) ON DELETE CASCADE,
    farmer_id UUID NOT NULL REFERENCES farmers(id) ON DELETE CASCADE,
    observation_type VARCHAR(50) NOT NULL,    -- 'fish_stress' | 'mortality' | 'discoloration' | 'algal_growth' | 'other'
    description TEXT NOT NULL,
    severity VARCHAR(20) NOT NULL,            -- 'low' | 'medium' | 'high'
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 7. Incidents Table
CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cage_id UUID NOT NULL REFERENCES cages(id) ON DELETE CASCADE,
    incident_type VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'open'         -- 'open' | 'resolved'
);

-- 8. Alerts Table
CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scope VARCHAR(20) NOT NULL,               -- 'cage' | 'farm' | 'region'
    related_id VARCHAR(255) NOT NULL,         -- cage_id, farm_id, or region_name
    severity VARCHAR(20) NOT NULL,            -- 'info' | 'warning' | 'critical'
    message TEXT NOT NULL,
    triggered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    acknowledged BOOLEAN DEFAULT FALSE
);

-- 9. Lake Zones Table (PostGIS Polygon for spatial lake regions)
CREATE TABLE IF NOT EXISTS lake_zones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    boundary GEOGRAPHY(Polygon, 4326) NOT NULL,
    region_label VARCHAR(100) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_lake_zones_boundary ON lake_zones USING GIST(boundary);

-- 10. Zone Metrics Table
CREATE TABLE IF NOT EXISTS zone_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    zone_id UUID NOT NULL REFERENCES lake_zones(id) ON DELETE CASCADE,
    period VARCHAR(20) NOT NULL,               -- e.g. '2026-08' or '10D-2026-08-20'
    avg_dissolved_oxygen NUMERIC(5, 2) NOT NULL,
    avg_temperature NUMERIC(5, 2) NOT NULL,
    avg_ph NUMERIC(4, 2) NOT NULL,
    avg_turbidity NUMERIC(6, 2) NOT NULL,
    risk_level VARCHAR(20) DEFAULT 'low',     -- 'low' | 'moderate' | 'high'
    trend VARCHAR(20) DEFAULT 'stable',       -- 'improving' | 'stable' | 'deteriorating'
    computed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 11. Expansion Signals Table
CREATE TABLE IF NOT EXISTS expansion_signals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    zone_id UUID NOT NULL REFERENCES lake_zones(id) ON DELETE CASCADE,
    suitability VARCHAR(50) NOT NULL,         -- 'high_suitability' | 'watch' | 'high_risk'
    rationale TEXT NOT NULL,
    computed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
