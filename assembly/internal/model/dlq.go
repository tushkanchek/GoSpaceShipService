package model

import "github.com/google/uuid"

type DeadLetterEvent struct {
	EventUuid       uuid.UUID              `json:"event_uuid"`
	FailedEventType string                 `json:"failed_event_type"`
	FailedEventData map[string]interface{} `json:"failed_event_data"`
	ErrorMessage    string                 `json:"error_message"`
	Timestamp       int64                  `json:"timestamp"`
}
