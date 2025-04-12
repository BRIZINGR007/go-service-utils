package helpers

import (
	"github.com/BRIZINGR007/go-service-utils/contextvars"
	"github.com/BRIZINGR007/go-service-utils/structs"
	"github.com/gin-gonic/gin"
)

func GetGinContextHeadersStruct(c *gin.Context) structs.Headers {
	contextvars := contextvars.ContextDataInit(c)
	return structs.Headers{
		CorrelationId: contextvars.GetCorrelationID(),
		Email:         contextvars.GetEmail(),
		UserId:        contextvars.GetUserID(),
		Authorization: contextvars.GetAuthToken(),
	}
}

func GetGinContextStringMap(c *gin.Context) map[string]string {
	ginContext := GetGinContextHeadersStruct(c)
	contextMap := map[string]string{
		"correlationId": ginContext.CorrelationId,
		"email":         ginContext.Email,
		"userId":        ginContext.UserId,
		"authorization": ginContext.Authorization,
	}
	return contextMap
}
