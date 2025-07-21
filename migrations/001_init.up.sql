CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_series (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id TEXT NOT NULL,
    PRIMARY KEY (user_id, series_id)
);
