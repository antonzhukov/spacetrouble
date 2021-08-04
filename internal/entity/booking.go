package entity

import "time"

// Booking describes details filled in the ticket
type Booking struct {
	ID            int64     `json:"id"`
	FirstName     string    `json:"first_name" validate:"required"`
	LastName      string    `json:"last_name" validate:"required"`
	Gender        string    `json:"gender"`
	Birthday      time.Time `json:"birthday" validate:"required"`
	LaunchpadID   string    `json:"launchpad_id" validate:"required"`
	DestinationID string    `json:"destination_id" validate:"required"`
	LaunchDate    time.Time `json:"launch_date" validate:"required"`
}
