-- migrate:up
CREATE TABLE
    IF NOT EXISTS bots (
        id BIGINT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        client_id VARCHAR(255) UNIQUE,
        username VARCHAR(255) NOT NULL,
        token BYTEA NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NULL
    );

-- Create index on client_id for fast lookups (already unique, but explicit index)
CREATE INDEX IF NOT EXISTS idx_bots_client_id ON bots (client_id);

CREATE TABLE
    IF NOT EXISTS bot_users (
        bot_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        -- User embedded fields
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NULL,
        username VARCHAR(255) NULL,
        photo_url TEXT NULL,
        is_premium BOOLEAN NULL,
        -- BotUser specific fields
        ip INET NOT NULL,
        user_agent TEXT NULL,
        language VARCHAR(10) NULL,
        last_login_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NULL,
        -- Primary key
        PRIMARY KEY (bot_id, user_id),
        -- Foreign key
        CONSTRAINT fk_bot_users_bot_id FOREIGN KEY (bot_id) REFERENCES bots (id) ON DELETE CASCADE
    );

-- Create index on bot_id for FK lookups
CREATE INDEX IF NOT EXISTS idx_bot_users_bot_id ON bot_users (bot_id);

-- migrate:down
DROP TABLE IF EXISTS bot_users;

-- Drop bots table
DROP TABLE IF EXISTS bots;