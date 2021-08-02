package launch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/antonzhukov/spacetrouble/internal/entity"
)

func Test_spaceX_GetLaunches(t *testing.T) {
	tests := []struct {
		name     string
		response string
		status   int
		want     entity.Launches
		wantErr  bool
	}{
		{
			"server failed",
			"",
			http.StatusInternalServerError,
			nil,
			true,
		},
		{
			"decode failed",
			"Our launches are better than yours",
			http.StatusOK,
			nil,
			true,
		},
		{
			"happy path",
			`[
    {
        "fairings": null,
        "launchpad": "5e9e4501f509094ba4566f84",
        "auto_update": true,
        "flight_number": 91,
        "name": "CRS-20",
        "date_utc": "2020-03-07T04:50:31.000Z",
        "upcoming": false
    }
]`,
			http.StatusOK,
			[]*entity.Launch{{
				Launchpad:    "5e9e4501f509094ba4566f84",
				AutoUpdate:   true,
				FlightNumber: 91,
				Name:         "CRS-20",
				DateUtc:      time.Date(2020, 03, 07, 04, 50, 31, 0, time.UTC),
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, tt.response)
			}))
			s := NewSpaceX(&http.Client{}, srv.URL)

			// act
			got, err := s.GetLaunches()
			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLaunches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLaunches() failed, %s", cmp.Diff(got, tt.want))
			}
		})
	}
}
