-- +goose Up
create table if not exists users (
    user_uuid varchar(36) primary key,
    login varchar(256) not null,
    password varchar(256) not null,
    email varchar(256),
    notifications jsonb,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- Создание триггерной функции
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER LANGUAGE plpgsql AS $$BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END$$;

-- Создание триггера
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS users;
