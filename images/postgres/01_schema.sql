CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
        id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		username VARCHAR(16) UNIQUE NOT NULL,
		password TEXT NOT NULL
);

CREATE TABLE posts (
        id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		post_owner VARCHAR(16) NOT NULL,
		content    TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (post_owner) REFERENCES users(username) ON DELETE CASCADE
);
