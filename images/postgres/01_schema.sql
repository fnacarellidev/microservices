CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
        id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		username VARCHAR(16) UNIQUE NOT NULL,
		password TEXT NOT NULL
);
