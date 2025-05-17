CREATE INDEX IF NOT EXISTS idx_parking_spot ON bookings(parking_spot);
CREATE INDEX IF NOT EXISTS idx_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_status ON bookings(status);