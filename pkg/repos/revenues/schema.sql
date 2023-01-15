CREATE MATERIALIZED VIEW IF NOT EXISTS by_book_revenues AS
  SELECT
    books.id,
    books.name,
    books.category,
    books.price,
    books.created_at,
    sum(orders.price) AS total,
    sum(
      CASE WHEN
        (orders.created_at > now() - interval '30 day')
      THEN orders.price ELSE 0 END
    ) AS last30d
  FROM
    books
    LEFT JOIN orders ON books.id = orders.book_id
  GROUP BY
      books.id;

CREATE UNIQUE INDEX IF NOT EXISTS v_books_id_unique_idx
  ON v_books (id);

CREATE UNIQUE INDEX IF NOT EXISTS v_books_total_volume_idx
  ON v_books (total);

CREATE UNIQUE INDEX IF NOT EXISTS v_books_total_last_30d_volume_idx
  ON v_books (last30d);
