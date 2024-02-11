CREATE TABLE tasks (
    id INT PRIMARY KEY,
    expression text NOT NULL,
    status text NOT NULL,
    created_at timestamp NOT NULL,
    calculated_at timestamp,
    calculated_by int
);