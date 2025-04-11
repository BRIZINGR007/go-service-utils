package structs

type MessageBody struct {
	Context map[string]string `json:"context" binding:"required"`
	Event   string            `json:"event" binding:"required"`
	Payload map[string]string `json:"payload" binding:"required"`
}
