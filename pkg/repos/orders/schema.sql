CREATE TABLE IF NOT EXISTS orders (
   id         INT         GENERATED ALWAYS AS IDENTITY,
   user_id    INT         references Users(ID) ON DELETE SET NULL,
   book_id    INT         references Items(ID) ON DELETE SET NULL,
   price      BIGINT      NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   is_deleted BOOLEAN     NOT NULL,
   CONSTRAINT orders_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS orders_item_id_idx ON orders (ItemID);
