package internal

import (
	"testing"
	"time"

	"github.com/antonzhukov/spacetrouble/internal/mocks"

	"github.com/google/go-cmp/cmp"

	"go.uber.org/zap"

	"github.com/antonzhukov/spacetrouble/internal/entity"
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
				&mocks.LaunchProvider{Launches: tt.launches},
				&mocks.BookingStore{AddErr: tt.bookingAddErr},
			)
			if err := bc.AddBooking(tt.b); (err != nil) != tt.wantErr {
				t.Errorf("AddBooking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
				&mocks.LaunchProvider{Launches: tt.launches},
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
			[]*entity.Booking{},
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
				{ID: 1, FirstName: "Anton", LaunchpadID: "launchpad-100"},
				{ID: 2, FirstName: "Viktor", LaunchpadID: "launchpad-101"},
				{ID: 3, FirstName: "Theodor", LaunchpadID: "launchpad-102"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mocks.BookingStore{}
			bc := NewBookingCenter(zap.L(), &mocks.LaunchProvider{}, store)
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
