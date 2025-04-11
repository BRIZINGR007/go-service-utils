package middlewares

import (
	"net/http"

	"github.com/BRIZINGR007/go-service-utils/contextvars"
	"github.com/BRIZINGR007/go-service-utils/helpers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RestMiddleware(context *gin.Context) {
	contextVars := contextvars.ContextDataInit(context)

	correlationId := context.Request.Header.Get("correlationId")
	if correlationId == "" {
		correlationId = uuid.New().String()
	}

	token := context.Request.Header.Get("authorization")
	if token == "" {
		cookieToken, err := context.Cookie("authorization")
		if err == nil {
			token = cookieToken
		}
	}
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"messsage": "Not Authorized"})
		return
	}
	claims, err := helpers.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	contextVars.SetUserID(claims.UserID)
	contextVars.SetEmail(claims.Email)
	contextVars.SetAuthToken(token)
	contextVars.SetCorrelationID(correlationId)
	context.Next()
}
