package handlers

import (
	"log"
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

		character, err := characterRepository.GetCharacterInDatabaseByName(ctx, req["name"])
		if err == nil {
			ctx.JSON(http.StatusOK, character)
			return
		}
		log.Println("error getting character in local database")

		character, err = characterRepository.GetCharacterInExternalAPIByName(ctx, req["name"])
		if err != nil {
			switch {
			case err == domains.ErrCharacterNotFoundInExternalAPI:
				ctx.JSON(http.StatusNotFound, gin.H{"error": err})
				return

			default:
				log.Println(err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}

		ctx.JSON(http.StatusOK, character)
	}
}
