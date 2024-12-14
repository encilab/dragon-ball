package handlers

import (
	"net/http"

	"github.com/encilab/dragon-ball/src/domains"
	"github.com/gin-gonic/gin"
)

func GetCharactersHandler(characterRepository domains.CharacterRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := make(map[string]string)
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": domains.ErrNameIsRequired})
			return
		}

		if req["name"] == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": domains.ErrNameIsRequired})
			return
		}

		character, err := characterRepository.GetCharacterInExternalAPIByName(ctx, req["name"])
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "character not found in external API"})
			return
		}

		ctx.JSON(http.StatusOK, character)
	}
}
