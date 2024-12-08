package models

// PayloadRequest defines the request format for executing a payload.
type PayloadRequest struct {
	PayloadName string   `json:"payload_name"`
	Options     []string `json:"options"`
}
