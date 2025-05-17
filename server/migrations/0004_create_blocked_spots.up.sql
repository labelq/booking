CREATE TABLE IF NOT EXISTS blocked_spots (
    spot_number INTEGER PRIMARY KEY,
    is_blocked BOOLEAN DEFAULT false,
    blocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Инициализируем все парковочные места как незаблокированные
INSERT INTO blocked_spots (spot_number, is_blocked)
SELECT generate_series(1, 16), false
ON CONFLICT (spot_number) DO NOTHING;