CREATE TABLE IF NOT EXISTS followers (
    pk_relation BIGSERIAL PRIMARY KEY,
    id_followee BIGINT NOT NULL,
    id_followed BIGINT NOT NULL,
    FOREIGN KEY (id_followee) REFERENCES users(pk_user),
    FOREIGN KEY (id_followed) REFERENCES users(pk_user)
);
