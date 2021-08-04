package internal

import (
	"time"

	"go.uber.org/zap"

	"github.com/antonzhukov/spacetrouble/internal/booking"
	"github.com/antonzhukov/spacetrouble/internal/entity"
	"github.com/antonzhukov/spacetrouble/internal/launch"

	"github.com/pkg/errors"
)

var FlightCancelledErr = errors.New("flight cancelled")

type BookingCenter struct {
	logger   *zap.Logger
	launches launch.Provider
	store    booking.Store
}

func NewBookingCenter(logger *zap.Logger, launches launch.Provider, store booking.Store) *BookingCenter {
	return &BookingCenter{
		logger:   logger,
		launches: launches,
		store:    store,
	}
}

// AddBooking adds a new booking in case flight operate criteria are met
func (bc *BookingCenter) AddBooking(b *entity.Booking) error {
	if b == nil {
		return errors.New("booking is empty")
	}

	if !bc.isLaunchPossible(b.LaunchDate, b.LaunchpadID) {
		return FlightCancelledErr
	}

	err := bc.store.Add(b)
	if err != nil {
		return errors.Wrap(err, "storing booking failed")
	}

	return nil
}

func (bc *BookingCenter) GetBookings() ([]*entity.Booking, error) {
	return bc.store.GetAll()
}

// isLaunchPossible finds if a flight on a given day can take place from a given launchpad
// going through existing competitor launches and checking if there's conflict
func (bc *BookingCenter) isLaunchPossible(date time.Time, launchpad string) bool {
	date = date.Truncate(24 * time.Hour)
	ll, err := bc.launches.GetLaunches()
	if err != nil {
		bc.logger.Info("GetLaunches failed", zap.Error(err))
		return false
	}
	for _, l := range ll {
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
