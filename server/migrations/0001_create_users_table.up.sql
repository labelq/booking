CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создаем администратора по умолчанию
INSERT INTO users (email, password_hash, account_type)
VALUES ('admin@example.com', '$2a$10$cAsJJCmKPe9HfBrN7srWIuvYOc9AZqGdleqzPcKNJd0sXVgIZPATO', 'admin')
    ON CONFLICT (email) DO NOTHING;