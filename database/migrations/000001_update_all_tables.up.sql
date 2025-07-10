-- Enums
CREATE TYPE document_type AS ENUM (
    'ID',
    'PASSPORT',
    'DRIVING_LICENSE',
    'NATIONAL_ID'
);

CREATE TYPE document_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TYPE payment_method AS ENUM ('CASH', 'MOBILE_MONEY', 'BANK_TRANSFER');

CREATE TYPE payment_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TYPE gender AS ENUM ('MALE', 'FEMALE', 'OTHER');

CREATE TYPE marital_status AS ENUM ('SINGLE', 'MARRIED', 'DIVORCED', 'WIDOWED');

CREATE TYPE religion AS ENUM (
    'CHRISTIAN',
    'MUSLIM',
    'HINDU',
    'BUDDHIST',
    'OTHER',
    'NONE'
);

-- Enable uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tables
CREATE TABLE region (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', name)) STORED
);

CREATE TABLE city (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    region_id UUID NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', name)) STORED,
    UNIQUE (name, region_id),
    FOREIGN KEY (region_id) REFERENCES region(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE address (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    house_number TEXT NOT NULL,
    district TEXT NOT NULL,
    city_id UUID NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', house_number || ' ' || district)
    ) STORED,
    UNIQUE (house_number, district, city_id),
    FOREIGN KEY (city_id) REFERENCES city(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE resident (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE CHECK (
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),
    first_name TEXT NOT NULL,
    second_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    birth_date TIMESTAMP(3) NOT NULL,
    gender gender NOT NULL,
    phone VARCHAR(20) NOT NULL CHECK (phone ~ '^\+?[1-9]\d{1,14}$'),
    marital_status marital_status NOT NULL,
    religion religion NOT NULL,
    ethnicity TEXT,
    disability_status TEXT,
    education_level TEXT,
    languages_spoken TEXT NOT NULL,
    address_id UUID,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            first_name || ' ' || second_name || ' ' || last_name
        )
    ) STORED,
    FOREIGN KEY (address_id) REFERENCES address(id) ON DELETE
    SET NULL ON UPDATE CASCADE
);

CREATE TABLE biometric (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id UUID NOT NULL UNIQUE,
    fingerprint BYTEA,
    blood_type TEXT NOT NULL,
    face TEXT NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', blood_type)) STORED,
    FOREIGN KEY (resident_id) REFERENCES resident(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE document (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    TYPE document_type NOT NULL,
    resident_id UUID NOT NULL UNIQUE,
    url TEXT NOT NULL,
    STATUS document_status NOT NULL,
    number TEXT NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', number)) STORED,
    FOREIGN KEY (resident_id) REFERENCES resident(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE idcard (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id UUID NOT NULL UNIQUE,
    number TEXT NOT NULL UNIQUE,
    issue_date TIMESTAMP(3) NOT NULL,
    expiry_date TIMESTAMP(3) NOT NULL,
    issue_place TEXT NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', number || ' ' || issue_place)
    ) STORED,
    FOREIGN KEY (resident_id) REFERENCES resident(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE payment (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id UUID NOT NULL UNIQUE,
    amount NUMERIC NOT NULL CHECK (amount >= 0),
    description TEXT NOT NULL,
    STATUS payment_status NOT NULL,
    reference TEXT NOT NULL,
    method payment_method NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', description || ' ' || reference)
    ) STORED,
    FOREIGN KEY (resident_id) REFERENCES resident(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE employment (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id UUID NOT NULL UNIQUE,
    STATUS TEXT NOT NULL,
    occupation TEXT,
    employer_name TEXT,
    work_address TEXT,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            STATUS || ' ' || COALESCE(occupation, '') || ' ' || COALESCE(employer_name, '') || ' ' || COALESCE(work_address, '')
        )
    ) STORED,
    FOREIGN KEY (resident_id) REFERENCES resident(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE emergency (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id UUID NOT NULL UNIQUE,
    name TEXT NOT NULL,
    relation TEXT NOT NULL,
    phone VARCHAR(20) NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', name || ' ' || relation)) STORED,
    FOREIGN KEY (resident_id) REFERENCES resident(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL,
    second_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE CHECK (
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),
    phone TEXT NOT NULL CHECK (phone ~ '^\+?[1-9]\d{1,14}$'),
    PASSWORD TEXT NOT NULL,
    role TEXT NOT NULL CHECK (
        role IN (
            'SUPERADMIN',
            'MANAGER',
            'ENCODER',
            'ADMIN',
            'CASHIER'
        )
    ),
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3)
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP(3) NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for full-text search
CREATE INDEX region_search_idx ON region USING GIN (search_vector);

CREATE INDEX city_search_idx ON city USING GIN (search_vector);

CREATE INDEX address_search_idx ON address USING GIN (search_vector);

CREATE INDEX resident_search_idx ON resident USING GIN (search_vector);

CREATE INDEX biometric_search_idx ON biometric USING GIN (search_vector);

CREATE INDEX document_search_idx ON document USING GIN (search_vector);

CREATE INDEX idcard_search_idx ON idcard USING GIN (search_vector);

CREATE INDEX payment_search_idx ON payment USING GIN (search_vector);

CREATE INDEX employment_search_idx ON employment USING GIN (search_vector);

CREATE INDEX emergency_search_idx ON emergency USING GIN (search_vector);