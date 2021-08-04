package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/antonzhukov/spacetrouble/internal"
	"github.com/antonzhukov/spacetrouble/internal/entity"
	"github.com/go-playground/validator/v10"
)

type Booking struct {
	logger *zap.Logger
	center *internal.BookingCenter
}

func NewBooking(logger *zap.Logger, center *internal.BookingCenter) *Booking {
	return &Booking{logger: logger, center: center}
}

// AddBooking handles new booking request
func (h Booking) AddBooking(w http.ResponseWriter, r *http.Request) {
	// decode JSON
	var b *entity.Booking
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		h.logger.Info("json.Decode failed", zap.Error(err))
		http.Error(w, "failed to decode json", http.StatusInternalServerError)
		return
	}

	// validate booking
	v := validator.New()
	err = v.Struct(b)
	if err != nil {
		h.logger.Info("validate booking failed", zap.Error(err))
		http.Error(w, "failed to validate booking", http.StatusBadRequest)
		return
	}

	// add booking
	err = h.center.AddBooking(b)
	if err != nil {
		h.logger.Info("AddBooking failed", zap.Error(err))
		if err == internal.FlightCancelledErr {
			http.Error(w, "flight cancelled", http.StatusBadRequest)
		} else {
			http.Error(w, "booking procedure failed", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetBookings return all available bookings
func (h Booking) GetBookings(w http.ResponseWriter, _ *http.Request) {
	bookings, err := h.center.GetBookings()
	if err != nil {
		h.logger.Info("GetBookings failed", zap.Error(err))
		http.Error(w, "get bookings failed", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(bookings)
	if err != nil {
		h.logger.Info("json.Encode failed", zap.Error(err))
		http.Error(w, "failed to encode bookings", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
