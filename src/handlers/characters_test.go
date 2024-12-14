package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/encilab/dragon-ball/src/domains"
	"github.com/encilab/dragon-ball/src/domains/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_GetCharactersHandler(t *testing.T) {
	name := "goku"
	characterDomain := domains.Character{
		ID:       1,
		Name:     "Goku",
		Ki:       60000000,
		Race:     "Saiyan",
		UrlImage: "https://dragonball-api.com/characters/goku_normal.webp",
	}

	characterRepoMock := mocks.NewCharacterRepository(t)
	characterRepoMock.On("GetCharacterInExternalAPIByName", mock.Anything, name).Return(characterDomain, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/characters", GetCharactersHandler(characterRepoMock))

	t.Run("given a valid request, it returns 200", func(t *testing.T) {

		testReq := map[string]string{
			"name": "goku",
		}

		testReqJSON, err := json.Marshal(testReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/characters", bytes.NewBuffer(testReqJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

}
