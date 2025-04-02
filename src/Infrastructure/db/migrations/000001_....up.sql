CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY,
    secret TEXT NOT NULL,
    redirect_uris TEXT[],
    grant_types TEXT[]
);
CREATE TABLE IF NOT EXISTS authorization_codes (
    code TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    client_id TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
