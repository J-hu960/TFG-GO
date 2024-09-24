CREATE TABLE IF NOT EXISTS reviews (
    pk_relation BIGSERIAL PRIMARY KEY,
    id_user BIGINT NOT NULL,
    id_project BIGINT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    content VARCHAR(500) not null,
    FOREIGN KEY (id_user) REFERENCES users(pk_user),
    FOREIGN KEY (id_project) REFERENCES projects(pk_project)
);
