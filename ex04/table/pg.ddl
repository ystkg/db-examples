CREATE TABLE shop (
  id serial PRIMARY KEY,
  name text NOT NULL UNIQUE,
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
