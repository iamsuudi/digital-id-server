-- Enums
CREATE TYPE document_type AS ENUM (
    'ID',
    'PASSPORT',
    'DRIVING_LICENSE',
    'NATIONAL_ID',
    'OTHER'
);
CREATE TYPE document_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');
CREATE TYPE payment_method AS ENUM ('CASH', 'MOBILE_MONEY', 'BANK_TRANSFER');
CREATE TYPE payment_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');
CREATE TYPE gender AS ENUM ('MALE', 'FEMALE', 'OTHER');
CREATE TYPE marital_status AS ENUM ('SINGLE', 'MARRIED', 'DIVORCED', 'WIDOWED');
CREATE TYPE religion AS ENUM (
    'CHRISTIAN',
    'MUSLIM',
    'ATHIEST',
    'HINDU',
    'BUDDHIST',
    'OTHER',
    'NONE'
);

-- Enable uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Tables
CREATE TABLE actor (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL,
    second_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE CHECK (
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),
    phone TEXT NOT NULL CHECK (phone ~ '^\+?[1-9]\d{1,14}$'),
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (
        role IN (
            'SUPERADMIN',
            'ADMIN',
            'MANAGER',
            'EXECUTIVE',
            'ENCODER',
            'CASHIER'
        )
    ),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            first_name || ' ' || second_name || ' ' || last_name
        )
    ) STORED,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3)
);

CREATE TABLE city (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    admin_id UUID,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', name)) STORED,
    UNIQUE (name),
    FOREIGN KEY (admin_id) REFERENCES actor(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE subcity (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    city_id UUID NOT NULL,
    manager_id UUID,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', name)) STORED,
    UNIQUE (name, city_id),
    FOREIGN KEY (city_id) REFERENCES city(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (manager_id) REFERENCES actor(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE kebele (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    subcity_id UUID,
    executive_id UUID,
    city_id UUID NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', name)) STORED,
    UNIQUE (name, city_id),
    FOREIGN KEY (city_id) REFERENCES city(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (subcity_id) REFERENCES subcity(id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (executive_id) REFERENCES actor(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE address (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    house_number TEXT NOT NULL,
    kebele_id UUID NOT NULL,
    subcity_id UUID,
    city_id UUID NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', house_number)
    ) STORED,
    UNIQUE (house_number, kebele_id, city_id),
    FOREIGN KEY (city_id) REFERENCES city(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (kebele_id) REFERENCES kebele(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (subcity_id) REFERENCES subcity(id) ON DELETE SET NULL ON UPDATE CASCADE
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
    status document_status NOT NULL,
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
    status payment_status NOT NULL,
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
    status TEXT NOT NULL,
    occupation TEXT,
    employer_name TEXT,
    work_address TEXT,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(3),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            status || ' ' || COALESCE(occupation, '') || ' ' || COALESCE(employer_name, '') || ' ' || COALESCE(work_address, '')
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

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP(3) NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES actor(id) ON DELETE CASCADE
);

-- Indexes for full-text search
CREATE INDEX actor_search_idx ON actor USING GIN (search_vector);
CREATE INDEX city_search_idx ON city USING GIN (search_vector);
CREATE INDEX subcity_search_idx ON city USING GIN (search_vector);
CREATE INDEX kebele_search_idx ON city USING GIN (search_vector);
CREATE INDEX address_search_idx ON address USING GIN (search_vector);
CREATE INDEX resident_search_idx ON resident USING GIN (search_vector);
CREATE INDEX biometric_search_idx ON biometric USING GIN (search_vector);
CREATE INDEX document_search_idx ON document USING GIN (search_vector);
CREATE INDEX idcard_search_idx ON idcard USING GIN (search_vector);
CREATE INDEX payment_search_idx ON payment USING GIN (search_vector);
CREATE INDEX employment_search_idx ON employment USING GIN (search_vector);
CREATE INDEX emergency_search_idx ON emergency USING GIN (search_vector);


-- Indexes for fuzzy search
-- 1. actor (full name)
CREATE INDEX actor_name_trgm_idx ON actor USING gin ((first_name || ' ' || second_name || ' ' || last_name) gin_trgm_ops);

-- 2. city
CREATE INDEX city_name_trgm_idx ON city USING gin (name gin_trgm_ops);

-- 3. subcity
CREATE INDEX subcity_name_trgm_idx ON subcity USING gin (name gin_trgm_ops);

-- 4. kebele
CREATE INDEX kebele_name_trgm_idx ON kebele USING gin (name gin_trgm_ops);

-- 5. address (house number)
CREATE INDEX address_house_trgm_idx ON address USING gin (house_number gin_trgm_ops);

-- 6. resident (full name)
CREATE INDEX resident_name_trgm_idx ON resident USING gin ((first_name || ' ' || second_name || ' ' || last_name) gin_trgm_ops);

-- 7. biometric (blood_type)
CREATE INDEX biometric_blood_trgm_idx ON biometric USING gin (blood_type gin_trgm_ops);

-- 8. document (number)
CREATE INDEX document_number_trgm_idx ON document USING gin (number gin_trgm_ops);

-- 9. idcard (number + issue_place)
CREATE INDEX idcard_combo_trgm_idx ON idcard USING gin ((number || ' ' || issue_place) gin_trgm_ops);

-- 10. payment (description + reference)
CREATE INDEX payment_combo_trgm_idx ON payment USING gin ((description || ' ' || reference) gin_trgm_ops);

-- 11. employment (status + occupation + employer_name + work_address)
CREATE INDEX employment_combo_trgm_idx ON employment USING gin (
    (status || ' ' || COALESCE(occupation,'') || ' ' ||
     COALESCE(employer_name,'') || ' ' || COALESCE(work_address,'')) gin_trgm_ops
);

-- 12. emergency (name + relation)
CREATE INDEX emergency_combo_trgm_idx ON emergency USING gin ((name || ' ' || relation) gin_trgm_ops);
