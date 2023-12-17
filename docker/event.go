package main

import "time"

const (
	// EventTypes
	GroupLifecycleCreate             = "group.lifecycle.create"
	GroupMembershipAdd               = "group.user_membership.add"
	GroupMembershipRemove            = "group.user_membership.remove"
	GroupProfileUpdate               = "group.profile.update"
	GroupApplicationAssignmentAdd    = "group.application_assignment.add"
	GroupApplicationAssignmentRemove = "group.application_assignment.remove"
)

type EventData struct {
	EventType string    `json:"eventType"`
	EventTime time.Time `json:"eventTime"`
	Data      struct {
		Events []Event `json:"events"`
	} `json:"data"`
}

type Event struct {
	UUID           string       `json:"uuid"`
	Published      time.Time    `json:"published"`
	EventType      string       `json:"eventType"`
	DisplayMessage string       `json:"displayMessage"`
	Severity       string       `json:"severity"`
	Client         Client       `json:"client"`
	Actor          Actor        `json:"actor"`
	Outcome        Outcome      `json:"outcome"`
	Target         []Target     `json:"target"`
	DebugContext   DebugContext `json:"debugContext"`
}

type Client struct {
	UserAgent           UserAgent           `json:"userAgent"`
	IPAddress           string              `json:"ipAddress"`
	GeographicalContext GeographicalContext `json:"geographicalContext"`
}

type UserAgent struct {
	RawUserAgent string `json:"rawUserAgent"`
	OS           string `json:"os"`
	Browser      string `json:"browser"`
}

type GeographicalContext struct {
	City        string      `json:"city"`
	State       string      `json:"state"`
	Country     string      `json:"country"`
	PostalCode  string      `json:"postalCode"`
	Geolocation Geolocation `json:"geolocation"`
}

type Geolocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Actor struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	AlternateId string `json:"alternateId"`
	DisplayName string `json:"displayName"`
}

type Outcome struct {
	Result string `json:"result"`
	Reason string `json:"reason,omitempty"`
}

type Target struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	AlternateId string `json:"alternateId"`
	DisplayName string `json:"displayName"`
}

type DebugContext struct {
	DebugData map[string]string `json:"debugData"`
}
