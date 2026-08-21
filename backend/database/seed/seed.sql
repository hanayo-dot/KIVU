-- Seed Script for KIVU Hackathon Demo (With Auth Credentials)
-- Pre-set Demo Farmer Credentials:
--   Phone: +254700000000
--   Password: demo1234
-- Bcrypt Hash for 'demo1234' (cost 12): $2a$12$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad6J1B4B48p.D9m

-- 1. Demo Farmer
INSERT INTO farmers (id, name, phone_number, password_hash, location_name)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Maina Farmer',
    '+254700000000',
    '$2a$12$BhtvxvGcqkMbLTcuoJPAf.VV6vs8bXyVXwiFnAdJTET4D7WtPEEPu',
    'Homa Bay Central'
) ON CONFLICT (phone_number) DO NOTHING;

-- 2. Demo Farm
INSERT INTO farms (id, farmer_id, name, region)
VALUES (
    'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Victoria Blue Aqua Farm',
    'Homa Bay Sector'
) ON CONFLICT (id) DO NOTHING;

-- 3. Demo Cages
INSERT INTO cages (id, farm_id, name, latitude, longitude, location, status)
VALUES 
(
    'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a31',
    'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'Cage A1 - Homa Bay North',
    -0.512300, 34.456700,
    ST_SetSRID(ST_MakePoint(34.456700, -0.512300), 4326)::geography,
    'active'
),
(
    'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a32',
    'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'Cage A2 - Homa Bay East',
    -0.515000, 34.460000,
    ST_SetSRID(ST_MakePoint(34.460000, -0.515000), 4326)::geography,
    'active'
),
(
    'c3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'Cage B1 - Mbita Channel',
    -0.430000, 34.200000,
    ST_SetSRID(ST_MakePoint(34.200000, -0.430000), 4326)::geography,
    'active'
),
(
    'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a34',
    'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'Cage C1 - Rusinga Deep',
    -0.400000, 34.150000,
    ST_SetSRID(ST_MakePoint(34.150000, -0.400000), 4326)::geography,
    'active'
) ON CONFLICT (id) DO NOTHING;

-- 4. Readings (Cage A1 showing deteriorating DO trend to trigger AI Anomaly Detection)
INSERT INTO readings (cage_id, dissolved_oxygen, temperature, ph, turbidity, recorded_at, source)
VALUES
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a31', 6.8, 25.4, 7.8, 12.0, NOW() - INTERVAL '3 days', 'sensor'),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a31', 5.9, 25.8, 7.7, 15.0, NOW() - INTERVAL '2 days', 'sensor'),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a31', 4.5, 26.5, 7.4, 22.0, NOW() - INTERVAL '1 day', 'sensor'),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a31', 2.8, 27.8, 7.1, 35.5, NOW() - INTERVAL '2 hours', 'sensor'),

('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a32', 6.5, 25.0, 7.6, 14.0, NOW() - INTERVAL '1 hour', 'sensor'),
('c3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 7.1, 24.8, 8.0, 9.5, NOW() - INTERVAL '30 minutes', 'sensor'),
('c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a34', 6.9, 25.1, 7.9, 10.0, NOW() - INTERVAL '15 minutes', 'manual');

-- 5. Observations
INSERT INTO observations (cage_id, farmer_id, observation_type, description, severity, recorded_at)
VALUES
(
    'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a31',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'fish_stress',
    'Fish surfacing near cage boundary gasping for air early morning.',
    'high',
    NOW() - INTERVAL '3 hours'
),
(
    'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a32',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'algal_growth',
    'Greenish tint floating on surface water.',
    'medium',
    NOW() - INTERVAL '1 day'
);

-- 6. Lake Zones (Simplified Polygons in Lake Victoria)
INSERT INTO lake_zones (id, name, boundary, region_label)
VALUES
(
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a41',
    'Homa Bay Sector',
    ST_SetSRID(ST_GeomFromText('POLYGON((34.40 -0.55, 34.50 -0.55, 34.50 -0.48, 34.40 -0.48, 34.40 -0.55))'), 4326)::geography,
    'Kenya Waters'
),
(
    'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a42',
    'Mbita & Rusinga Channel',
    ST_SetSRID(ST_GeomFromText('POLYGON((34.10 -0.45, 34.30 -0.45, 34.30 -0.35, 34.10 -0.35, 34.10 -0.45))'), 4326)::geography,
    'Kenya Waters'
),
(
    'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a43',
    'Winam Gulf Offshore',
    ST_SetSRID(ST_GeomFromText('POLYGON((34.50 -0.40, 34.75 -0.40, 34.75 -0.15, 34.50 -0.15, 34.50 -0.40))'), 4326)::geography,
    'Kenya Waters'
),
(
    'b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
    'Mwanza Gulf Region',
    ST_SetSRID(ST_GeomFromText('POLYGON((32.80 -2.60, 33.10 -2.60, 33.10 -2.30, 32.80 -2.30, 32.80 -2.60))'), 4326)::geography,
    'Tanzania Waters'
) ON CONFLICT (id) DO NOTHING;

-- 7. Zone Metrics (Satellite Composites)
INSERT INTO zone_metrics (zone_id, period, avg_dissolved_oxygen, avg_temperature, avg_ph, avg_turbidity, risk_level, trend)
VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a41', '2026-08', 3.2, 27.5, 7.2, 38.0, 'high', 'deteriorating'),
('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a42', '2026-08', 6.8, 25.1, 7.9, 11.2, 'low', 'improving'),
('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a43', '2026-08', 5.4, 26.0, 7.6, 21.0, 'moderate', 'stable'),
('b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', '2026-08', 7.1, 24.5, 8.1, 8.5, 'low', 'stable');

-- 8. Expansion Signals
INSERT INTO expansion_signals (zone_id, suitability, rationale)
VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a41', 'high_risk', 'High turbidity and low dissolved oxygen levels detected via Sentinel-3 LWQ composite.'),
('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a42', 'high_suitability', 'Optimal temperature (25.1°C), low turbidity (11.2 NTU), and strong natural water exchange.'),
('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a43', 'watch', 'Moderate nutrient loading; monitoring seasonal turn-over.'),
('b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'high_suitability', 'Excellent dissolved oxygen stability and clean open water.');

-- 9. Pre-existing Alerts
INSERT INTO alerts (scope, related_id, severity, message, triggered_at, acknowledged)
VALUES
('cage', 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a31', 'critical', 'Critical Hypoxia Alert: Dissolved oxygen dropped to 2.8 mg/L in Cage A1.', NOW() - INTERVAL '2 hours', false),
('farm', 'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'warning', 'FARM WARNING: Hypoxia detected across Homa Bay Sector cages.', NOW() - INTERVAL '1 hour', false),
('region', 'Homa Bay Sector', 'warning', 'REGIONAL WATCH: Elevated turbidity composite reported across Homa Bay Sector.', NOW() - INTERVAL '30 minutes', false);
