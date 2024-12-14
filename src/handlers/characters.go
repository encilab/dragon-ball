package handlers

import (
	"log"
	"net/http"
	"strconv"

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

func SearchCharactersHandler(characterRepository domains.CharacterRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var limit int
		if ctx.Query("limit") == "" {
			limit = 100
		} else {
			limitConvert, err := strconv.Atoi(ctx.Query("limit"))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			limit = limitConvert
		}

		character, err := characterRepository.SearchCharactersInDatabase(ctx, limit)
		if err != nil {
			switch {
			default:
				log.Println(err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}

		ctx.JSON(http.StatusOK, character)
	}
}

func DeleteCharacterHandler(characterRepository domains.CharacterRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")

		err := characterRepository.DeleteCharacterInDatabase(ctx, name)
		if err != nil {
			switch {
			case err == domains.ErrCharacterNotDeleted:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			default:
				log.Println(err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}

		ctx.Status(http.StatusOK)
	}
}
