-- Password reset tokens table
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP(3) NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    used_at TIMESTAMP(3),
    CONSTRAINT valid_token CHECK (expires_at > created_at)
);

-- Create indexes for performance
CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);

CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);

CREATE INDEX idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);
