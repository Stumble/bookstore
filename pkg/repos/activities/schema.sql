CREATE TABLE IF NOT EXISTS activities (
   id            INT                 GENERATED ALWAYS AS IDENTITY,
   action        VARCHAR(255)        NOT NULL,
   parameter     TEXT,
   created_at    TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
   CONSTRAINT activities_id_pkey PRIMARY KEY (id)
);
