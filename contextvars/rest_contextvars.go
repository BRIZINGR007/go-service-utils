package contextvars

import (
	"github.com/gin-gonic/gin"
)

type ContextData struct {
	Ctx *gin.Context
}

func ContextDataInit(ctx *gin.Context) *ContextData {
	return &ContextData{Ctx: ctx}
}

func (h *ContextData) SetUserID(userID string) {
	h.Ctx.Set("userId", userID)
}

func (h *ContextData) GetUserID() string {
	val, exists := h.Ctx.Get("userId")
	if !exists {
		return ""
	}
	return val.(string)
}

func (h *ContextData) SetEmail(email string) {
	h.Ctx.Set("email", email)
}

func (h *ContextData) GetEmail() string {
	val, exists := h.Ctx.Get("email")
	if !exists {
		return ""
	}
	return val.(string)
}

func (h *ContextData) SetCorrelationID(id string) {
	h.Ctx.Set("correlationId", id)
}

func (h *ContextData) GetCorrelationID() string {
	val, exists := h.Ctx.Get("correlationId")
	if !exists {
		return ""
	}
	return val.(string)
}

func (h *ContextData) SetAuthToken(token string) {
	h.Ctx.Set("authorization", token)
}

func (h *ContextData) GetAuthToken() string {
	val, exists := h.Ctx.Get("authorization")
	if !exists {
		return ""
	}
	return val.(string)
}
