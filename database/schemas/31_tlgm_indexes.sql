-- Trigram GIN indexes for LIKE / ILIKE / fuzzy search
CREATE INDEX IF NOT EXISTS user_name_trgm_idx       ON "user"     USING GIN ((first_name||' '||second_name||' '||last_name) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS city_name_trgm_idx       ON city       USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS subcity_name_trgm_idx    ON subcity    USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS kebele_name_trgm_idx     ON kebele     USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS address_house_trgm_idx   ON address    USING GIN (house_number gin_trgm_ops);
CREATE INDEX IF NOT EXISTS resident_name_trgm_idx   ON resident   USING GIN ((first_name||' '||second_name||' '||last_name) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS biometric_blood_trgm_idx ON biometric  USING GIN (blood_type gin_trgm_ops);
CREATE INDEX IF NOT EXISTS document_number_trgm_idx ON document   USING GIN (number gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idcard_combo_trgm_idx    ON idcard     USING GIN ((number||' '||issue_place) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS payment_combo_trgm_idx   ON payment    USING GIN ((description||' '||reference) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS employment_combo_trgm_idx ON employment USING GIN (
  (status||' '||COALESCE(occupation,'')||' '||
   COALESCE(employer_name,'')||' '||COALESCE(work_address,'')) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS emergency_combo_trgm_idx ON emergency  USING GIN ((name||' '||relation) gin_trgm_ops);