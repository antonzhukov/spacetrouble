package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/antonzhukov/spacetrouble/internal/entity"

	"github.com/google/go-cmp/cmp"

	"github.com/antonzhukov/spacetrouble/internal/mocks"

	"github.com/antonzhukov/spacetrouble/internal"
	"go.uber.org/zap"
)

func TestBooking(t *testing.T) {
	tests := []struct {
		name       string
		booking    string
		want       string
		wantStatus int
	}{
		{
			"booking added",
			`{"first_name":"Anton", "last_name":"Zhukov", "birthday":"2021-01-01T00:00:00+00:00", "launchpad_id":"5e9e4502f509092b78566f87", "destination_id": "abc", "launch_date": "2031-01-01T00:00:00+00:00"}`,
			`[{"id":1,"first_name":"Anton","last_name":"Zhukov","gender":"","birthday":"2021-01-01T00:00:00Z","launchpad_id":"5e9e4502f509092b78566f87","destination_id":"abc","launch_date":"2031-01-01T00:00:00Z"}]
`,
			http.StatusOK,
		},
		{
			"flight cancelled",
			`{"first_name":"Anton", "last_name":"Zhukov", "birthday":"2021-01-01T00:00:00+00:00", "launchpad_id":"5e9e4502f509092b78566f87", "destination_id": "abc", "launch_date": "2022-01-01T00:00:00+00:00"}`,
			`[]
`,
			http.StatusBadRequest,
		},
		{
			"bad json",
			`abc`,
			`[]
`,
			http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			launches := &mocks.LaunchProvider{
				Launches: []*entity.Launch{{
					Launchpad: "5e9e4502f509092b78566f87",
					DateUtc:   time.Date(2022, 01, 01, 15, 0, 0, 0, time.UTC),
					Upcoming:  true,
				}},
			}
			bc := internal.NewBookingCenter(zap.L(), launches, &mocks.BookingStore{})
			h := NewBooking(zap.L(), bc)
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(tt.booking)))
			h.AddBooking(resp, req)

			if resp.Code != tt.wantStatus {
				t.Errorf("AddBooking failed, %s", cmp.Diff(resp.Code, tt.wantStatus))
			}

			resp = httptest.NewRecorder()
			req, _ = http.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(tt.booking)))
			h.GetBookings(resp, req)
			got := resp.Body.String()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("AddBooking failed, %s", cmp.Diff(got, tt.want))
			}
		})
	}
}
