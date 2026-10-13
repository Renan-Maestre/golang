-- This is a sample migration.

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  biography TEXT NOT NULL
  
)

---- create above / drop below ----

DROP TABLE IF EXISTS users;
