-- GIN indexes for pre-computed search_vector columns
CREATE INDEX IF NOT EXISTS user_search_idx      ON "user"      USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS city_search_idx       ON city       USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS subcity_search_idx    ON subcity    USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS kebele_search_idx     ON kebele     USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS address_search_idx    ON address    USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS resident_search_idx   ON resident   USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS biometric_search_idx  ON biometric  USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS document_search_idx   ON document   USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS idcard_search_idx     ON idcard     USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS payment_search_idx    ON payment    USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS employment_search_idx ON employment USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS emergency_search_idx  ON emergency  USING GIN (search_vector);
