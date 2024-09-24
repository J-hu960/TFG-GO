CREATE TABLE IF NOT EXISTS user_disliked_projects (
    pk_relation BIGSERIAL PRIMARY KEY,
    id_user BIGINT NOT NULL,
    id_project BIGINT NOT NULL,
    FOREIGN KEY (id_user) REFERENCES users(pk_user),
    FOREIGN KEY (id_project) REFERENCES projects(pk_project)
);
