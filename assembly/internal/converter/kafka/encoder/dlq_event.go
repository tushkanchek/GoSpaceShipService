package encoder

import (
	"encoding/json"
	"fmt"

	"assembly/internal/model"
)

type dlqEncoder struct{}

func NewDLQEventEncoder() *dlqEncoder {
	return &dlqEncoder{}
}

func (e *dlqEncoder) Encode(event model.DeadLetterEvent) ([]byte, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("json marshal dlq event: %w", err)
	}
	return data, nil
}
