CREATE TABLE computing_resources (
    id SERIAL PRIMARY KEY,
    name text,
    occupied bool DEFAULT FALSE
);