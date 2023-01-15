CREATE TABLE IF NOT EXISTS users (
   id          INT          GENERATED ALWAYS AS IDENTITY,
   name        VARCHAR(255) NOT NULL,
   metadata    JSON,
   image       TEXT         NOT NULL,
   created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
   CONSTRAINT users_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS users_created_at_idx
    ON Users (CreatedAt);

CREATE UNIQUE INDEX IF NOT EXISTS users_lower_name_key
    ON Users ((lower(Name))) INCLUDE (ID);
