CREATE TABLE bookings (
                          id SERIAL PRIMARY KEY,
                          user_id INTEGER NOT NULL,
                          parking_spot INTEGER NOT NULL CHECK (parking_spot > 0 AND parking_spot <= 16),
                          car_number VARCHAR(20) NOT NULL,
                          reserved_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          hours INTEGER NOT NULL CHECK (hours > 0),
                          FOREIGN KEY (user_id) REFERENCES users(id)
);
