package design

import (
	"fmt"
	"sync"
	"time"
)

/**
 * Design Movie Ticket Booking System
 */

type Movie struct {
    ID       string
    Title    string
    Duration time.Duration
    Genre    string
}

type Show struct {
    ID        string
    Movie     *Movie
    StartTime time.Time
    TheaterID string
    ScreenID  string
}

type Seat struct {
    Row      int
    Col      int
    IsBooked bool
}

type Screen struct {
    ID    string
    Seats [][]*Seat
}

type Theater struct {
    ID      string
    Name    string
    Screens map[string]*Screen
}

type Booking struct {
    ID        string
    ShowID    string
    Seats     []*Seat
    UserID    string
    Status    string // Pending, Confirmed, Cancelled
    CreatedAt time.Time
}

type BookingService struct {
    mu       sync.Mutex
    bookings map[string]*Booking
}

func generateID() string {
    return ""
}

func (bs *BookingService) BookSeats(show *Show, seats []*Seat, userID string) (*Booking, error) {
    bs.mu.Lock()
    defer bs.mu.Unlock()

    // check availability
    for _, seat := range seats {
        if seat.IsBooked {
            return nil, fmt.Errorf("seat already booked")
        }
    }

    // mark seats
    for _, seat := range seats {
        seat.IsBooked = true
    }

    booking := &Booking{
        ID:        generateID(),
        ShowID:    show.ID,
        Seats:     seats,
        UserID:    userID,
        Status:    "Confirmed",
        CreatedAt: time.Now(),
    }
    bs.bookings[booking.ID] = booking

    return booking, nil
}
