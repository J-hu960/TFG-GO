CREATE TYPE user_role AS ENUM ('user', 'admin', 'moderator');


CREATE TABLE IF NOT EXISTS users(
    pk_user bigserial PRIMARY KEY,
    -- phone  VARCHAR(50) NOT NULL UNIQUE default 'default',
    created_at timestamp(0) with time zone NOT NULL default NOW(),
    email VARCHAR(255) NOT NULL UNIQUE,
    profile_pict TEXT DEFAULT 'default' not null,
    hashed_password bytea not null,
    description text default 'none' not null,
    role user_role DEFAULT 'user' not null
);

CREATE INDEX idx_users_email ON users(email);