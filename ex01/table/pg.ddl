CREATE TABLE movie (
  id serial PRIMARY KEY,
  title text NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL
);
