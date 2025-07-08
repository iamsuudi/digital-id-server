-- Drop indexes first (reverse order of creation)
DROP INDEX IF EXISTS "emergency_search_idx";

DROP INDEX IF EXISTS "employment_search_idx";

DROP INDEX IF EXISTS "payment_search_idx";

DROP INDEX IF EXISTS "idcard_search_idx";

DROP INDEX IF EXISTS "document_search_idx";

DROP INDEX IF EXISTS "biometric_search_idx";

DROP INDEX IF EXISTS "resident_search_idx";

DROP INDEX IF EXISTS "address_search_idx";

DROP INDEX IF EXISTS "city_search_idx";

DROP INDEX IF EXISTS "region_search_idx";

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS "emergency";

DROP TABLE IF EXISTS "employment";

DROP TABLE IF EXISTS "payment";

DROP TABLE IF EXISTS "idcard";

DROP TABLE IF EXISTS "document";

DROP TABLE IF EXISTS "biometric";

DROP TABLE IF EXISTS "resident";

DROP TABLE IF EXISTS "address";

DROP TABLE IF EXISTS "city";

DROP TABLE IF EXISTS "region";

-- Drop enums last
DROP TYPE IF EXISTS "Religion";

DROP TYPE IF EXISTS "MaritalStatus";

DROP TYPE IF EXISTS "Gender";

DROP TYPE IF EXISTS "PaymentStatus";

DROP TYPE IF EXISTS "PaymentMethod";

DROP TYPE IF EXISTS "DocumentStatus";

DROP TYPE IF EXISTS "DocumentType";