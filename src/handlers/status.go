package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LivezHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		responseMap := make(map[string]string)
		responseMap["message"] = "ok"

		ctx.JSON(http.StatusOK, responseMap)
	}
}
