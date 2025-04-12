package structs

import "encoding/json"

type MessageBody struct {
	Context map[string]string `json:"context" binding:"required"`
	Event   string            `json:"event" binding:"required"`
	Payload json.RawMessage   `json:"payload" binding:"required"`
}
