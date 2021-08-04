package mocks

import (
	"github.com/antonzhukov/spacetrouble/internal/entity"
	"github.com/pkg/errors"
)

type LaunchProvider struct {
	Launches []*entity.Launch
}

func (m *LaunchProvider) GetLaunches() (entity.Launches, error) {
	return m.Launches, nil
}

type BookingStore struct {
	bookings []*entity.Booking
	AddErr   bool
}

func (m *BookingStore) Add(b *entity.Booking) error {
	if m.AddErr {
		return errors.New("failed")
	}
	b.ID = int64(len(m.bookings) + 1)
	m.bookings = append(m.bookings, b)
	return nil
}

func (m *BookingStore) GetAll() ([]*entity.Booking, error) {
	if len(m.bookings) == 0 {
		return []*entity.Booking{}, nil
	}
	return m.bookings, nil
}
