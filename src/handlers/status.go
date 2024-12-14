package handlers

import (
	"database/sql"
	"log"
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

func ReadyzHandler(sqlClient *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := sqlClient.Ping(); err != nil {
			log.Println("error when execute sqlClient.Ping, err: " + err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database is down"})
			return
		}

		responseMap := make(map[string]string)
		responseMap["message"] = "ok"

		ctx.JSON(http.StatusOK, responseMap)
	}
}
