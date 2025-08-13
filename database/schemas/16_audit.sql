CREATE TABLE audit_log (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_id            UUID NOT NULL REFERENCES "user"(id) ON DELETE SET NULL,
    target_user_id      UUID REFERENCES "user"(id) ON DELETE SET NULL,
    target_resident_id  UUID REFERENCES resident(id) ON DELETE SET NULL,
    target_role_slug    TEXT,
    action_type         TEXT NOT NULL,
    object_type         TEXT NOT NULL,
    object_id           BIGINT,
    diff                JSONB,
    ts                  TIMESTAMP(3) NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_actor     ON audit_log(actor_id);
CREATE INDEX idx_audit_target    ON audit_log(target_user_id);
CREATE INDEX idx_audit_role      ON audit_log(target_role_slug);
CREATE INDEX idx_audit_ts        ON audit_log(ts DESC);
CREATE INDEX idx_audit_object    ON audit_log(object_type, object_id);
