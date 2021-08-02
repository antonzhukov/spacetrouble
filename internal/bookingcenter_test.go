package internal

import (
	"testing"
	"time"

	"github.com/antonzhukov/spacetrouble/internal/entity"
	"github.com/pkg/errors"
)

func TestBookingsCenter_AddBooking(t *testing.T) {
	tests := []struct {
		name          string
		launches      []*entity.Launch
		b             *entity.Booking
		bookingAddErr bool
		wantErr       bool
	}{
		{
			"booking details empty",
			nil,
			nil,
			false,
			true,
		},
		{
			"launch not possible",
			[]*entity.Launch{{
				Upcoming:  true,
				DateUtc:   time.Date(2060, 01, 01, 13, 0, 0, 0, time.UTC),
				Launchpad: "launchpad-100",
			}},
			&entity.Booking{
				LaunchpadID: "launchpad-100",
				LaunchDate:  time.Date(2060, 01, 01, 13, 0, 0, 0, time.UTC),
			},
			false,
			true,
		},
		{
			"booking add error",
			[]*entity.Launch{{
				Upcoming:  true,
				DateUtc:   time.Date(2060, 01, 01, 13, 0, 0, 0, time.UTC),
				Launchpad: "launchpad-100",
			}},
			&entity.Booking{
				LaunchpadID: "launchpad-101",
				LaunchDate:  time.Date(2060, 01, 01, 13, 0, 0, 0, time.UTC),
			},
			true,
			true,
		},
		{
			"happy path",
			[]*entity.Launch{{
				Upcoming:  true,
				DateUtc:   time.Date(2060, 01, 01, 13, 0, 0, 0, time.UTC),
				Launchpad: "launchpad-100",
			}},
			&entity.Booking{
				LaunchpadID: "launchpad-101",
				LaunchDate:  time.Date(2060, 01, 01, 13, 0, 0, 0, time.UTC),
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BookingCenter{
				launches: &mockLaunchProvider{launches: tt.launches},
				store:    &mockBookingStore{addErr: tt.bookingAddErr},
			}
			if err := m.AddBooking(tt.b); (err != nil) != tt.wantErr {
				t.Errorf("AddBooking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockLaunchProvider struct {
	launches []*entity.Launch
}

func (m *mockLaunchProvider) GetLaunches() []*entity.Launch {
	return m.launches
}

type mockBookingStore struct {
	addErr bool
}

func (m *mockBookingStore) Add(b *entity.Booking) error {
	if m.addErr {
		return errors.New("failed")
	}
	return nil
}

func TestBookingsCenter_IsLaunchPossible(t *testing.T) {
	tests := []struct {
		name      string
		launches  []*entity.Launch
		date      time.Time
		launchpad string
		want      bool
	}{
		{
			"no actual launches, launch possible",
			[]*entity.Launch{nil, {DateUtc: time.Time{}}},
			time.Date(2060, 01, 01, 0, 0, 0, 0, time.UTC),
			"launchpad-100",
			true,
		},
		{
			"conflicting launch, launch not possible",
			[]*entity.Launch{{
				Upcoming:  true,
				DateUtc:   time.Date(2060, 01, 01, 13, 0, 0, 0, time.UTC),
				Launchpad: "launchpad-100",
			}},
			time.Date(2060, 01, 01, 0, 0, 0, 0, time.UTC),
			"launchpad-100",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BookingCenter{
				launches: &mockLaunchProvider{launches: tt.launches},
			}
			if got := m.isLaunchPossible(tt.date, tt.launchpad); got != tt.want {
				t.Errorf("isLaunchPossible() = %v, want %v", got, tt.want)
			}
		})
	}
}
