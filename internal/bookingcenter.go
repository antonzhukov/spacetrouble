package internal

import (
	"fmt"
	"time"

	"github.com/antonzhukov/spacetrouble/internal/booking"
	"github.com/antonzhukov/spacetrouble/internal/entity"
	"github.com/antonzhukov/spacetrouble/internal/launch"

	"github.com/pkg/errors"
)

type BookingCenter struct {
	launches launch.Provider
	store    booking.Store
}

// AddBooking adds a new booking in case flight operate criteria are met
func (m *BookingCenter) AddBooking(b *entity.Booking) error {
	if b == nil {
		return errors.New("booking is empty")
	}

	if !m.isLaunchPossible(b.LaunchDate, b.LaunchpadID) {
		return fmt.Errorf("launch is cancelled at %s", b.LaunchDate.String())
	}

	err := m.store.Add(b)
	if err != nil {
		return errors.Wrap(err, "storing booking failed")
	}

	return nil
}

// isLaunchPossible finds if a flight on a given day can take place from a given launchpad
// going through existing competitor launches and checking if there's conflict
func (m *BookingCenter) isLaunchPossible(date time.Time, launchpad string) bool {
	date = date.Truncate(24 * time.Hour)
	for _, l := range m.launches.GetLaunches() {
		if l == nil {
			continue
		}
		// only need launches in the future
		if !l.Upcoming {
			continue
		}

		if l.Launchpad == launchpad && l.DateUtc.Truncate(24*time.Hour).Equal(date) {
			return false
		}
	}

	return true
}
