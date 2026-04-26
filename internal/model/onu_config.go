package model

import "time"

// ONUExecuteRequest is the generic payload for ONU write operations.
type ONUExecuteRequest struct {
	OLTID     string         `json:"olt_id"`
	Operation string         `json:"operation"`
	Target    string         `json:"target"`
	DryRun    bool           `json:"dry_run"`
	OnuRef    ONUReference   `json:"onu_ref"`
	Payload   map[string]any `json:"payload"`
	Confirm   bool           `json:"confirm"`
}

// ONUReference identifies ONU location.
type ONUReference struct {
	Slot  int `json:"slot"`
	PON   int `json:"pon"`
	OnuID int `json:"onu_id"`
}

// ONUExecuteResult represents dry-run or execute outcome.
type ONUExecuteResult struct {
	Mode      string         `json:"mode"`
	Message   string         `json:"message"`
	Applied   bool           `json:"applied"`
	Plan      []string       `json:"operation_plan,omitempty"`
	Warnings  []string       `json:"warnings,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	Meta      map[string]any `json:"meta,omitempty"`
}

// OLTProfile stores one OLT connection profile.
type OLTProfile struct {
	ID        string    `json:"olt_id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      uint16    `json:"port"`
	Community string    `json:"community,omitempty"`
	TimeoutMS int       `json:"timeout_ms"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OLTProfileRequest is create/update request body.
type OLTProfileRequest struct {
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      uint16 `json:"port"`
	Community string `json:"community"`
	TimeoutMS int    `json:"timeout_ms"`
}
