package entity

import "time"

// Booking describes details filled in the ticket
type Booking struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Gender        string    `json:"gender"`
	Birthday      time.Time `json:"birthday"`
	LaunchpadID   string    `json:"launchpad_id"`
	DestinationID string    `json:"destination_id"`
	LaunchDate    time.Time `json:"launch_date"`
}
