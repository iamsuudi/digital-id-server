CREATE TABLE resident (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email               TEXT NOT NULL UNIQUE,
    first_name          TEXT NOT NULL,
    second_name         TEXT NOT NULL,
    last_name           TEXT NOT NULL,
    birth_date          TIMESTAMP(3)    NOT NULL,
    gender              VARCHAR(20)     NOT NULL,
    phone               VARCHAR(20)     NOT NULL,
    marital_status      VARCHAR(20),
    religion            VARCHAR(20),
    national_id         INTEGER UNIQUE,
    ethnicity           VARCHAR(20),
    disability          VARCHAR(30),
    education_level     VARCHAR(20),
    languages_spoken    VARCHAR(20)[],
    address_id          UUID REFERENCES address(id) ON DELETE SET NULL,
    
    created_at          TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP(3)
);
