CREATE TABLE setting (
    id                          TEXT PRIMARY KEY,
    idcard_expiration_duration  INT NOT NULL,
    created_at                  TIMESTAMP(3) DEFAULT NOW(),
    updated_at                  TIMESTAMP(3) DEFAULT NOW()
);

INSERT INTO setting (id, idcard_expiration_duration) VALUES ('settings', 5);
