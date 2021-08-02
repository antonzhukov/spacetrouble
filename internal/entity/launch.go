package entity

import "time"

type Launches []*Launch

// Launch describes the entity as defined in
// https://github.com/r-spacex/SpaceX-API/blob/master/docs/launches/v4/all.md
type Launch struct {
	Fairings interface{} `json:"fairings"`
	Links    struct {
		Patch struct {
			Small string `json:"small"`
			Large string `json:"large"`
		} `json:"patch"`
		Reddit struct {
			Campaign string      `json:"campaign"`
			Launch   string      `json:"launch"`
			Media    string      `json:"media"`
			Recovery interface{} `json:"recovery"`
		} `json:"reddit"`
		Flickr struct {
			Small    []string `json:"small"`
			Original []string `json:"original"`
		} `json:"flickr"`
		Presskit  string `json:"presskit"`
		Webcast   string `json:"webcast"`
		YoutubeId string `json:"youtube_id"`
		Article   string `json:"article"`
		Wikipedia string `json:"wikipedia"`
	} `json:"links"`
	StaticFireDateUtc  time.Time     `json:"static_fire_date_utc"`
	StaticFireDateUnix int           `json:"static_fire_date_unix"`
	Tdb                bool          `json:"tdb"`
	Net                bool          `json:"net"`
	Window             int           `json:"window"`
	Rocket             string        `json:"rocket"`
	Success            bool          `json:"success"`
	Failures           []interface{} `json:"failures"`
	Details            string        `json:"details"`
	Crew               []interface{} `json:"crew"`
	Ships              []interface{} `json:"ships"`
	Capsules           []string      `json:"capsules"`
	Payloads           []string      `json:"payloads"`
	Launchpad          string        `json:"launchpad"`
	AutoUpdate         bool          `json:"auto_update"`
	FlightNumber       int           `json:"flight_number"`
	Name               string        `json:"name"`
	DateUtc            time.Time     `json:"date_utc"`
	DateUnix           int           `json:"date_unix"`
	DateLocal          time.Time     `json:"date_local"`
	DatePrecision      string        `json:"date_precision"`
	Upcoming           bool          `json:"upcoming"`
	Cores              []struct {
		Core           string `json:"core"`
		Flight         int    `json:"flight"`
		Gridfins       bool   `json:"gridfins"`
		Legs           bool   `json:"legs"`
		Reused         bool   `json:"reused"`
		LandingAttempt bool   `json:"landing_attempt"`
		LandingSuccess bool   `json:"landing_success"`
		LandingType    string `json:"landing_type"`
		Landpad        string `json:"landpad"`
	} `json:"cores"`
	Id string `json:"id"`
}
