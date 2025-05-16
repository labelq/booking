-- Создаем расширение для UUID если оно еще не создано
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаем последовательности для ID
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE SEQUENCE IF NOT EXISTS bookings_id_seq;
CREATE SEQUENCE IF NOT EXISTS blocked_spots_id_seq;

-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY DEFAULT nextval('users_id_seq'),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создаем таблицу бронирований
CREATE TABLE IF NOT EXISTS bookings (
                                        id INTEGER PRIMARY KEY DEFAULT nextval('bookings_id_seq'),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    parking_spot INTEGER NOT NULL,
    car_number VARCHAR(20) NOT NULL,
    reserved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    hours INTEGER NOT NULL CHECK (hours > 0),
    CONSTRAINT valid_parking_spot CHECK (parking_spot > 0 AND parking_spot <= 16)
    );

-- Создаем таблицу заблокированных мест
CREATE TABLE IF NOT EXISTS blocked_spots (
                                             id INTEGER PRIMARY KEY DEFAULT nextval('blocked_spots_id_seq'),
    spot_number INTEGER NOT NULL,
    is_blocked BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_spot_number UNIQUE(spot_number),
    CONSTRAINT valid_spot_number CHECK (spot_number > 0 AND spot_number <= 16)
    );

-- Создаем индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_parking_spot ON bookings(parking_spot);
CREATE INDEX IF NOT EXISTS idx_blocked_spots_number ON blocked_spots(spot_number);

-- Создаем триггерную функцию для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

-- Создаем триггер для обновления updated_at в таблице blocked_spots
DROP TRIGGER IF EXISTS update_blocked_spots_updated_at ON blocked_spots;
CREATE TRIGGER update_blocked_spots_updated_at
    BEFORE UPDATE ON blocked_spots
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Инициализируем все парковочные места как незаблокированные
INSERT INTO blocked_spots (spot_number, is_blocked)
SELECT generate_series(1, 16), false
    ON CONFLICT (spot_number) DO NOTHING;

-- Предоставляем необходимые права (если нужно)
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO postgres;