CREATE TABLE IF NOT EXISTS operations (
    id SERIAL PRIMARY KEY,
    name text NOT NULL UNIQUE,
    duration int DEFAULT 5
);