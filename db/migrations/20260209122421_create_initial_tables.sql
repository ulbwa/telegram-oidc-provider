-- migrate:up
-- Create users table
CREATE TABLE
    users (
        id BIGINT PRIMARY KEY,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255),
        username VARCHAR(255),
        photo_url TEXT,
        is_premium BOOLEAN,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    );

-- Create bots table
CREATE TABLE
    bots (
        id BIGINT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        client_id VARCHAR(255) NOT NULL,
        username VARCHAR(255) NOT NULL,
        token VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    );

-- Create unique index on bots.client_id for fast lookup by client_id
CREATE UNIQUE INDEX idx_bots_client_id ON bots (client_id);

-- Create user_bot_logins junction table
CREATE TABLE
    user_bot_logins (
        user_id BIGINT NOT NULL,
        bot_id BIGINT NOT NULL,
        ip INET NOT NULL,
        user_agent TEXT,
        language VARCHAR(10),
        last_login_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP,
        PRIMARY KEY (user_id, bot_id),
        CONSTRAINT fk_user_bot_logins_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
        CONSTRAINT fk_user_bot_logins_bot_id FOREIGN KEY (bot_id) REFERENCES bots (id) ON DELETE CASCADE
    );

-- Create index on bot_id for efficient queries by bot
CREATE INDEX idx_user_bot_logins_bot_id ON user_bot_logins (bot_id);

-- Create composite index for efficient pagination by bot_id with last_login_at ordering
CREATE INDEX idx_user_bot_logins_bot_id_last_login ON user_bot_logins (bot_id, last_login_at DESC);

-- migrate:down
-- Drop tables in reverse order
DROP TABLE IF EXISTS user_bot_logins;

DROP TABLE IF EXISTS bots;

DROP TABLE IF EXISTS users;