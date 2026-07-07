package design

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_BookingService(t *testing.T) {
	show := &Show{
		ID:        "show1",
		Movie:     &Movie{ID: "m1", Title: "Inception", Duration: 2 * time.Hour, Genre: "Sci-Fi"},
		StartTime: time.Now(),
		TheaterID: "t1",
		ScreenID:  "s1",
	}

	t.Run("successful_booking", func(t *testing.T) {
		bs := &BookingService{bookings: map[string]*Booking{}}
		seats := []*Seat{{Row: 1, Col: 1}, {Row: 1, Col: 2}}

		booking, err := bs.BookSeats(show, seats, "user1")
		assert.NoError(t, err)
		assert.NotNil(t, booking)
		assert.Equal(t, "show1", booking.ShowID)
		assert.Equal(t, "user1", booking.UserID)
		assert.Equal(t, "Confirmed", booking.Status)
		assert.True(t, seats[0].IsBooked)
		assert.True(t, seats[1].IsBooked)
	})

	t.Run("double_booking_returns_error", func(t *testing.T) {
		bs := &BookingService{bookings: map[string]*Booking{}}
		seat := &Seat{Row: 2, Col: 3, IsBooked: true}

		_, err := bs.BookSeats(show, []*Seat{seat}, "user2")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already booked")
	})

	t.Run("partial_conflict_rejects_all", func(t *testing.T) {
		bs := &BookingService{bookings: map[string]*Booking{}}
		free := &Seat{Row: 3, Col: 1}
		taken := &Seat{Row: 3, Col: 2, IsBooked: true}

		_, err := bs.BookSeats(show, []*Seat{free, taken}, "user3")
		assert.Error(t, err)
		assert.False(t, free.IsBooked) // not marked because error returned early
	})
}
