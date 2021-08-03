package internal

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"go.uber.org/zap"

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
			bc := NewBookingCenter(
				zap.L(),
				&mockLaunchProvider{launches: tt.launches},
				&mockBookingStore{addErr: tt.bookingAddErr},
			)
			if err := bc.AddBooking(tt.b); (err != nil) != tt.wantErr {
				t.Errorf("AddBooking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockLaunchProvider struct {
	launches []*entity.Launch
}

func (m *mockLaunchProvider) GetLaunches() (entity.Launches, error) {
	return m.launches, nil
}

type mockBookingStore struct {
	bookings []*entity.Booking
	addErr   bool
}

func (m *mockBookingStore) Add(b *entity.Booking) error {
	if m.addErr {
		return errors.New("failed")
	}
	m.bookings = append(m.bookings, b)
	return nil
}

func (m *mockBookingStore) GetAll() ([]*entity.Booking, error) {
	return m.bookings, nil
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
			bc := NewBookingCenter(
				zap.L(),
				&mockLaunchProvider{launches: tt.launches},
				nil,
			)
			if got := bc.isLaunchPossible(tt.date, tt.launchpad); got != tt.want {
				t.Errorf("isLaunchPossible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookingCenter_GetBookings(t *testing.T) {
	tests := []struct {
		name     string
		bookings []*entity.Booking
		want     []*entity.Booking
		wantErr  bool
	}{
		{
			"empty",
			nil,
			nil,
			false,
		},
		{
			"with values",
			[]*entity.Booking{
				{FirstName: "Anton", LaunchpadID: "launchpad-100"},
				{FirstName: "Viktor", LaunchpadID: "launchpad-101"},
				{FirstName: "Theodor", LaunchpadID: "launchpad-102"},
			},
			[]*entity.Booking{
				{FirstName: "Anton", LaunchpadID: "launchpad-100"},
				{FirstName: "Viktor", LaunchpadID: "launchpad-101"},
				{FirstName: "Theodor", LaunchpadID: "launchpad-102"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockBookingStore{}
			bc := NewBookingCenter(zap.L(), &mockLaunchProvider{}, store)
			for _, b := range tt.bookings {
				bc.AddBooking(b)
			}

			got, err := bc.GetBookings()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBookings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetBookings() got failed, %s", cmp.Diff(got, tt.want))
			}
		})
	}
}
