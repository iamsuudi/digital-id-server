-- name: GetGenderDistribution :many
SELECT
    gender,
    COUNT(*) AS count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 1) AS percentage
FROM resident_analytics
GROUP BY gender
ORDER BY count DESC;

-- name: GetAgeGroupDistribution :many
SELECT
    age_group,
    COUNT(*) AS count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 1) AS percentage
FROM resident_analytics
GROUP BY age_group
ORDER BY
    CASE age_group
        WHEN 'Under 18' THEN 1
        WHEN '18-24' THEN 2
        WHEN '25-34' THEN 3
        WHEN '35-44' THEN 4
        WHEN '45-54' THEN 5
        WHEN '55-64' THEN 6
        ELSE 7
    END;

-- name: GetGenderAgeGroupDistribution :many
SELECT
    gender,
    age_group,
    COUNT(*) AS count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY gender), 1) AS percentage_within_gender
FROM resident_analytics
GROUP BY gender, age_group
ORDER BY gender,
    CASE age_group
        WHEN 'Under 18' THEN 1
        WHEN '18-24' THEN 2
        WHEN '25-34' THEN 3
        WHEN '35-44' THEN 4
        WHEN '45-54' THEN 5
        WHEN '55-64' THEN 6
        ELSE 7
    END;

-- name: GetRegistrationTrends :many
SELECT
    DATE_TRUNC('month', created_at) AS month,
    COUNT(*) AS new_residents
FROM resident_analytics
GROUP BY month
ORDER BY month;
