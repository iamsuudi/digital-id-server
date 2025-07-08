-- Enums
CREATE TYPE "DocumentType" AS ENUM (
    'ID',
    'PASSPORT',
    'DRIVING_LICENSE',
    'NATIONAL_ID'
);

CREATE TYPE "DocumentStatus" AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TYPE "PaymentMethod" AS ENUM ('CASH', 'MOBILE_MONEY', 'BANK_TRANSFER');

CREATE TYPE "PaymentStatus" AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TYPE "Gender" AS ENUM ('MALE', 'FEMALE', 'OTHER');

CREATE TYPE "MaritalStatus" AS ENUM ('SINGLE', 'MARRIED', 'DIVORCED', 'WIDOWED');

CREATE TYPE "Religion" AS ENUM (
    'CHRISTIAN',
    'MUSLIM',
    'HINDU',
    'BUDDHIST',
    'OTHER',
    'NONE'
);

-- Tables
CREATE TABLE "region" (
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL UNIQUE,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', "name")) STORED
);

CREATE TABLE "city" (
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "regionId" SERIAL NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', "name")) STORED,
    UNIQUE ("name", "regionId"),
    FOREIGN KEY ("regionId") REFERENCES "region"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "address" (
    "id" SERIAL PRIMARY KEY,
    "houseNumber" TEXT NOT NULL,
    "district" TEXT NOT NULL,
    "cityId" SERIAL NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', "houseNumber" || ' ' || "district")
    ) STORED,
    UNIQUE ("houseNumber", "district", "cityId"),
    FOREIGN KEY ("cityId") REFERENCES "city"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "resident" (
    "id" SERIAL PRIMARY KEY,
    "email" TEXT NOT NULL UNIQUE,
    -- CHECK ("email" ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    "firstName" TEXT NOT NULL,
    "secondName" TEXT NOT NULL,
    "lastName" TEXT NOT NULL,
    "birthDate" TIMESTAMP(3) NOT NULL,
    "gender" "Gender" NOT NULL,
    "phone" VARCHAR(20) NOT NULL,
    "maritalStatus" "MaritalStatus" NOT NULL,
    "religion" "Religion" NOT NULL,
    "ethnicity" TEXT,
    "disabilityStatus" TEXT,
    "educationLevel" TEXT,
    "languagesSpoken" TEXT NOT NULL,
    "addressId" SERIAL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            "firstName" || ' ' || "secondName" || ' ' || "lastName"
        )
    ) STORED,
    FOREIGN KEY ("addressId") REFERENCES "address"("id") ON DELETE
    SET NULL ON UPDATE CASCADE
);

CREATE TABLE "biometric" (
    "id" SERIAL PRIMARY KEY,
    "residentId" SERIAL NOT NULL UNIQUE,
    "fingerprint" BYTEA,
    "bloodType" TEXT NOT NULL,
    "face" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', "bloodType")) STORED,
    FOREIGN KEY ("residentId") REFERENCES "resident"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "document" (
    "id" SERIAL PRIMARY KEY,
    "type" "DocumentType" NOT NULL,
    "residentId" SERIAL NOT NULL UNIQUE,
    "url" TEXT NOT NULL,
    "status" "DocumentStatus" NOT NULL,
    "number" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', "number")) STORED,
    FOREIGN KEY ("residentId") REFERENCES "resident"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "idcard" (
    "id" SERIAL PRIMARY KEY,
    "residentId" SERIAL NOT NULL UNIQUE,
    "number" TEXT NOT NULL UNIQUE,
    "issueDate" TIMESTAMP(3) NOT NULL,
    "expiryDate" TIMESTAMP(3) NOT NULL,
    "issuePlace" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', "number" || ' ' || "issuePlace")
    ) STORED,
    FOREIGN KEY ("residentId") REFERENCES "resident"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "payment" (
    "id" SERIAL PRIMARY KEY,
    "residentId" SERIAL NOT NULL UNIQUE,
    "amount" DOUBLE PRECISION NOT NULL CHECK ("amount" >= 0),
    "description" TEXT NOT NULL,
    "status" "PaymentStatus" NOT NULL,
    "reference" TEXT NOT NULL,
    "method" "PaymentMethod" NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', "description" || ' ' || "reference")
    ) STORED,
    FOREIGN KEY ("residentId") REFERENCES "resident"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "employment" (
    "id" SERIAL PRIMARY KEY,
    "residentId" SERIAL NOT NULL UNIQUE,
    "status" TEXT NOT NULL,
    "occupation" TEXT,
    "employerName" TEXT,
    "workAddress" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            "status" || ' ' || COALESCE("occupation", '') || ' ' || COALESCE("employerName", '') || ' ' || COALESCE("workAddress", '')
        )
    ) STORED,
    FOREIGN KEY ("residentId") REFERENCES "resident"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "emergency" (
    "id" SERIAL PRIMARY KEY,
    "residentId" SERIAL NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    "relation" TEXT NOT NULL,
    "phone" VARCHAR(20) NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),
    "searchVector" TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', "name" || ' ' || "relation")
    ) STORED,
    FOREIGN KEY ("residentId") REFERENCES "resident"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- Indexes for full-text search
CREATE INDEX "region_search_idx" ON "region" USING GIN ("searchVector");

CREATE INDEX "city_search_idx" ON "city" USING GIN ("searchVector");

CREATE INDEX "address_search_idx" ON "address" USING GIN ("searchVector");

CREATE INDEX "resident_search_idx" ON "resident" USING GIN ("searchVector");

CREATE INDEX "biometric_search_idx" ON "biometric" USING GIN ("searchVector");

CREATE INDEX "document_search_idx" ON "document" USING GIN ("searchVector");

CREATE INDEX "idcard_search_idx" ON "idcard" USING GIN ("searchVector");

CREATE INDEX "payment_search_idx" ON "payment" USING GIN ("searchVector");

CREATE INDEX "employment_search_idx" ON "employment" USING GIN ("searchVector");

CREATE INDEX "emergency_search_idx" ON "emergency" USING GIN ("searchVector");