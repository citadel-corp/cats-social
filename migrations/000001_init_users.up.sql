CREATE TABLE IF NOT EXISTS
users (
    id SERIAL PRIMARY KEY,
    uid CHAR(16) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    hashed_password BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS users_email
	ON users USING HASH (email);

CREATE INDEX IF NOT EXISTS users_uid
	ON users USING HASH (uid);
