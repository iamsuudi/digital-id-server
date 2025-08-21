CREATE OR REPLACE VIEW resident_analytics AS
SELECT
    id,
    email,
    first_name,
    second_name,
    last_name,
    birth_date,
    gender,
    phone,
    address_id,
    created_at,
    deleted_at,
    -- Calculate age in years
    EXTRACT(YEAR FROM AGE(birth_date)) AS age,
    -- Define age groups
    CASE
        WHEN EXTRACT(YEAR FROM AGE(birth_date)) < 18 THEN 'Under 18'
        WHEN EXTRACT(YEAR FROM AGE(birth_date)) BETWEEN 18 AND 24 THEN '18-24'
        WHEN EXTRACT(YEAR FROM AGE(birth_date)) BETWEEN 25 AND 34 THEN '25-34'
        WHEN EXTRACT(YEAR FROM AGE(birth_date)) BETWEEN 35 AND 44 THEN '35-44'
        WHEN EXTRACT(YEAR FROM AGE(birth_date)) BETWEEN 45 AND 54 THEN '45-54'
        WHEN EXTRACT(YEAR FROM AGE(birth_date)) BETWEEN 55 AND 64 THEN '55-64'
        ELSE '65+'
    END AS age_group
FROM resident
WHERE deleted_at IS NULL;  -- Exclude soft-deleted records;

-- Create indixes
CREATE INDEX idx_resident_gender ON resident(gender);
CREATE INDEX idx_resident_birth_date ON resident(birth_date);
CREATE INDEX idx_resident_created_at ON resident(created_at);
