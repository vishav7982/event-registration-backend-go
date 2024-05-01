package middlewares

import (
	"net/http"
	"restapp/utils"

	"github.com/gin-gonic/gin"
)

func Autheticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		//No token passed in the headers
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization token is not passed in the headers"})
		return
	}
	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	context.Set("userId", userId)
}
