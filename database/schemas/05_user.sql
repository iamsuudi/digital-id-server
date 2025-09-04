-- 1. roles (tree) ------------------------------------------------
CREATE TABLE role (
    slug                TEXT PRIMARY KEY,
    name                TEXT NOT NULL,
    parent_role_slug    TEXT REFERENCES role(slug) ON DELETE RESTRICT,
    level_rank          int  NOT NULL   -- lower = more powerful
);

-- 2. user -------------------------------------------------
CREATE TABLE "user" (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name       TEXT NOT NULL,
    second_name      TEXT NOT NULL,
    last_name        TEXT NOT NULL,
    email            TEXT UNIQUE NOT NULL,
    phone            TEXT NOT NULL,
    picture          TEXT,
    city_id          UUID REFERENCES city(id)    ON DELETE SET NULL,
    subcity_id       UUID REFERENCES subcity(id) ON DELETE SET NULL,
    kebele_id        UUID REFERENCES kebele(id)  ON DELETE SET NULL,
    role_slug        TEXT NOT NULL REFERENCES role(slug)  ON DELETE RESTRICT,

    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP(3)
);

CREATE TABLE account (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    password_hash   TEXT NOT NULL,
    user_id         UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    created_at      TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP(3)
);

-- 4. permissions -------------------------------------------------
CREATE TABLE permission (
    name        TEXT PRIMARY KEY,
    label       TEXT NOT NULL,
    description TEXT
);

INSERT INTO role (slug, name, parent_role_slug, level_rank) VALUES
('superadmin', 'SuperAdmin',  NULL, 0),
('admin',      'Admin',       'superadmin', 10),
('manager',    'Manager',     'admin', 20),
('executive',  'Executive',   'manager', 30),
('cashier',    'Cashier',     'executive', 40),
('encoder',    'Encoder',     'executive', 40);

-- Add a few permissions
INSERT INTO permission (name, label, description) VALUES
('can_edit_price',     'Edit Price',     'Modify product prices'),
('can_issue_refund',   'Issue Refund',   'Refund a sale'),
('can_void_receipt',   'Void Receipt',   'Cancel a finished receipt'),
('can_close_shift',    'Close Shift',    'End cashiering session'),
('can_export_reports', 'Export Reports', 'Download analytics');

-- 5. role_permissions (permissions belong only to roles) ---------
CREATE TABLE role_permission (
    role_slug       TEXT  REFERENCES role(slug) ON DELETE CASCADE,
    permission_name TEXT  REFERENCES permission(name) ON DELETE CASCADE,
    PRIMARY KEY (role_slug, permission_name)
);

-- Grant superadmin all current permissions
INSERT INTO role_permission (role_slug, permission_name) VALUES
('superadmin', 'can_edit_price'),
('superadmin', 'can_issue_refund'),
('superadmin', 'can_void_receipt'),
('superadmin', 'can_close_shift'),
('superadmin', 'can_export_reports')
ON CONFLICT DO NOTHING;

-- 7. sparse overrides (per-user permission toggles) --------------
CREATE TABLE user_permission_override (
    user_id         UUID         NOT NULL REFERENCES "user"(id)       ON DELETE CASCADE,
    permission_name TEXT         NOT NULL REFERENCES permission(name) ON DELETE CASCADE,
    is_granted      BOOLEAN      NOT NULL,
    granted_by      UUID         NOT NULL REFERENCES "user"(id)       ON DELETE SET NULL,
    granted_at      TIMESTAMP(3) NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, permission_name)
);

-- 8. useful indexes ---------------------------------------------
CREATE INDEX idx_role_parent            ON role(parent_role_slug);
CREATE INDEX idx_role_level             ON role(level_rank);
CREATE INDEX idx_user_scope             ON "user"(city_id, subcity_id, kebele_id);
CREATE INDEX idx_user_role              ON "user"(role_slug);
CREATE INDEX idx_role_perm_role         ON role_permission(role_slug);
CREATE INDEX idx_override_user          ON user_permission_override(user_id);
