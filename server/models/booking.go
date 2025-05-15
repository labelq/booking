package models

import (
    "database/sql"
    "time"
)

type Booking struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    ParkingSpot int       `json:"parking_spot"`
    CarNumber   string    `json:"car_number"`
    ReservedAt  time.Time `json:"reserved_at"`
    Hours       int       `json:"hours"`
}

func CreateBooking(db *sql.DB, booking *Booking) (int, error) {
    var id int
    err := db.QueryRow(`
        INSERT INTO bookings (user_id, parking_spot, car_number, reserved_at, hours)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, booking.UserID, booking.ParkingSpot, booking.CarNumber, booking.ReservedAt, booking.Hours).Scan(&id)

    return id, err
}

func GetOccupiedParkingSpots(db *sql.DB) ([]int, error) {
    spots := []int{}
    rows, err := db.Query(`
        SELECT DISTINCT parking_spot
        FROM bookings
        WHERE reserved_at + (hours * interval '1 hour') > NOW()
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var spot int
        if err := rows.Scan(&spot); err != nil {
            return nil, err
        }
        spots = append(spots, spot)
    }

    return spots, nil
}

func IsParkingSpotOccupied(db *sql.DB, spotNumber int) (bool, error) {
    var count int
    err := db.QueryRow(`
        SELECT COUNT(*)
        FROM bookings
        WHERE parking_spot = $1
        AND reserved_at + (hours * interval '1 hour') > NOW()
    `, spotNumber).Scan(&count)

    return count > 0, err
}