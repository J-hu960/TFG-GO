CREATE TYPE PROJECT_CATEGORY AS ENUM('sport', 'comedy', 'environment', 'technology');

CREATE TABLE IF NOT EXISTS projects (
    pk_project BIGSERIAL PRIMARY KEY,
    name VARCHAR(250) NOT NULL UNIQUE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    photos TEXT[] NOT NULL,
    link_web VARCHAR(300) NOT NULL DEFAULT 'default',
    description VARCHAR(500) NOT NULL,
    likes INT NOT NULL DEFAULT 0,
    dislikes INT NOT NULL DEFAULT 0,
    founds_recieved INT NOT NULL DEFAULT 0,
    founds_expected INT NOT NULL,
    category PROJECT_CATEGORY[] NOT NULL,  -- Si es un array
    id_creator BIGINT NOT NULL REFERENCES users(pk_user) -- Corrigiendo la clave foránea
);
