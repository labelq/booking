CREATE TABLE IF NOT EXISTS bookings (
                          id SERIAL PRIMARY KEY,
                          user_id INT NOT NULL,
                          parking_spot INT NOT NULL,
                          reserved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          FOREIGN KEY (user_id) REFERENCES users(id)
);