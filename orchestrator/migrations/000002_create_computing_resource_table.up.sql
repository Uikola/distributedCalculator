CREATE TABLE computing_resources (
    id SERIAL PRIMARY KEY,
    name text,
    task text,
    occupied bool DEFAULT FALSE
);