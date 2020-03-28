-- migrate:up
CREATE TABLE users (
    id SERIAL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_id ON users(id);

-- migrate:down

