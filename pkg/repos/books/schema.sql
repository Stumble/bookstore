CREATE TYPE category AS ENUM (
    'computer_science',
    'philosophy',
    'comic'
);

CREATE TABLE IF NOT EXISTS books (
   id            BIGSERIAL           GENERATED ALWAYS AS IDENTITY,
   name          VARCHAR(255)        NOT NULL,
   description   VARCHAR(255)        NOT NULL,
   metadata      JSON,
   category      ItemCategory        NOT NULL,
   price         DECIMAL(10,2)       NOT NULL,
   created_at    TIMESTAMP           NOT NULL DEFAULT NOW(),
   updated_at    TIMESTAMP           NOT NULL DEFAULT NOW(),
   CONSTRAINT books_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS books_name_idx ON books (name);
CREATE INDEX IF NOT EXISTS books_category_id_idx ON books (category, id);
