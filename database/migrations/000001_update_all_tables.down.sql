-- Drop indexes (must be dropped before tables)
DROP INDEX IF EXISTS emergency_search_idx;
DROP INDEX IF EXISTS employment_search_idx;
DROP INDEX IF EXISTS payment_search_idx;
DROP INDEX IF EXISTS idcard_search_idx;
DROP INDEX IF EXISTS document_search_idx;
DROP INDEX IF EXISTS biometric_search_idx;
DROP INDEX IF EXISTS resident_search_idx;
DROP INDEX IF EXISTS address_search_idx;
DROP INDEX IF EXISTS city_search_idx;
DROP INDEX IF EXISTS region_search_idx;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS emergency;
DROP TABLE IF EXISTS employment;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS idcard;
DROP TABLE IF EXISTS document;
DROP TABLE IF EXISTS biometric;
DROP TABLE IF EXISTS resident;
DROP TABLE IF EXISTS address;
DROP TABLE IF EXISTS city;
DROP TABLE IF EXISTS region;

-- Drop enums
DROP TYPE IF EXISTS religion;
DROP TYPE IF EXISTS marital_status;
DROP TYPE IF EXISTS gender;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS document_status;
DROP TYPE IF EXISTS document_type;
DROP EXTENSION IF EXISTS "uuid-ossp";