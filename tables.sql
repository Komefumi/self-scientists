CREATE TABLE rule_sets (
  identifier VARCHAR(255) PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp,
  can_delete_any_user_account boolean DEFAULT false,
  can_delete_any_thread boolean DEFAULT false,
  can_delete_any_post boolean DEFAULT false
);

CREATE TABLE site_account_types (
  identifier VARCHAR(255) PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp,
  rule_set_name VARCHAR(255) REFERENCES rule_sets (identifier) ON DELETE SET NULL
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  date_of_birth DATE NOT NULL,
  password_hash TEXT,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp,
  site_account_type VARCHAR(255) REFERENCES site_account_types (identifier) ON DELETE SET NULL,
  bio TEXT
);

CREATE TABLE categories (
  identifier VARCHAR(255) PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp
);

CREATE TABLE threads (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp,
  category_identifier VARCHAR(255) REFERENCES categories (identifier) ON DELETE SET NULL
);

CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  thread_id integer NOT NULL REFERENCES threads (id) ON DELETE CASCADE,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp,
  author_id integer NOT NULL REFERENCES users (id) ON DELETE SET NULL,
  replying_to_id integer REFERENCES posts (id) ON DELETE SET NULL
);
