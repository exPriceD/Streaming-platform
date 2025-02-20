CREATE TABLE users
(
    id                         UUID PRIMARY KEY                  DEFAULT gen_random_uuid(), -- Генерация UUID автоматически
    username                   VARCHAR(50)              NOT NULL UNIQUE,                    -- Уникальное имя пользователя, длина ограничена
    email                      VARCHAR(255)             NOT NULL UNIQUE,                    -- Уникальный email, длина может быть больше
    password_hash              TEXT                     NOT NULL,                           -- Хэш пароля (может быть длинным, поэтому TEXT)
    avatar_url                 TEXT,                                                        -- URL аватара (может быть длинным)
    consent_to_data_processing BOOLEAN                  NOT NULL DEFAULT FALSE,             -- Согласие на обработку данных
    created_at                 TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Время создания
    updated_at                 TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP  -- Время обновления
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_username ON users (username);