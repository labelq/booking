CREATE TABLE IF NOT EXISTS blocked_spots (
                                             spot_number INTEGER PRIMARY KEY,
                                             is_blocked BOOLEAN DEFAULT true,
                                             blocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);