CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- заголовок
    title TEXT NOT NULL,
    -- текст объявления
    card_text TEXT NOT NULL,
    -- адрес изображения
    picture_url TEXT NOT NULL,
    -- цена
    price NUMERIC NOT NULL,
    -- автор
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- дата создания
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);