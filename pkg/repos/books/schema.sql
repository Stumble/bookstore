CREATE TYPE book_category AS ENUM (
    'computer_science',
    'philosophy',
    'comic'
);

CREATE TABLE IF NOT EXISTS books (
   id            SERIAL              NOT NULL,
   name          VARCHAR(255)        NOT NULL,
   description   VARCHAR(255)        NOT NULL,
   metadata      JSON,
   category      book_category       NOT NULL,
   price         REAL                NOT NULL,
   dummy_field   INT,
   created_at    TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
   updated_at    TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
   CONSTRAINT books_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS books_name_idx ON books (name);
CREATE INDEX IF NOT EXISTS books_category_id_idx ON books (category, id);
