package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
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
		ID:    1,
		Name:  "Goku",
		Ki:    "60.000.000",
		Race:  "Saiyan",
		Image: "https://dragonball-api.com/characters/goku_normal.webp",
	}

	t.Run("given a valid request, it returns 200 when getting data of external api", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("GetCharacterInDatabaseByName", mock.Anything, name).Return(characterDomain, domains.ErrCharacterNotFoundInDatabase)
		characterRepoMock.On("GetCharacterInExternalAPIByName", mock.Anything, name).Return(characterDomain, nil)

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.POST("/api/characters", GetCharactersHandler(characterRepoMock))

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

	t.Run("given a valid request, it returns 200 when getting data of database", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("GetCharacterInDatabaseByName", mock.Anything, name).Return(characterDomain, nil)

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.POST("/api/characters", GetCharactersHandler(characterRepoMock))

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

	t.Run("given a valid request, it returns 400 when not send json in body data", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.POST("/api/characters", GetCharactersHandler(characterRepoMock))

		req, err := http.NewRequest(http.MethodPost, "/api/characters", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a valid request, it returns 400 when send name empty in body request", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.POST("/api/characters", GetCharactersHandler(characterRepoMock))

		testReq := map[string]string{
			"name": "",
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

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a valid request, it returns 500 when character not found in external api", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("GetCharacterInDatabaseByName", mock.Anything, name).Return(characterDomain, domains.ErrCharacterNotFoundInDatabase)
		characterRepoMock.On("GetCharacterInExternalAPIByName", mock.Anything, name).Return(characterDomain, domains.ErrCharacterNotFoundInExternalAPI)

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.POST("/api/characters", GetCharactersHandler(characterRepoMock))

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

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("given a valid request, it returns 500 when unexpected error", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("GetCharacterInDatabaseByName", mock.Anything, name).Return(characterDomain, domains.ErrCharacterNotFoundInDatabase)
		characterRepoMock.On("GetCharacterInExternalAPIByName", mock.Anything, name).Return(characterDomain, errors.New("any error"))

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.POST("/api/characters", GetCharactersHandler(characterRepoMock))

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

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func Test_SearchCharactersHandler(t *testing.T) {
	limit := 100
	characterDomain := domains.Character{
		ID:    1,
		Name:  "Goku",
		Ki:    "60.000.000",
		Race:  "Saiyan",
		Image: "https://dragonball-api.com/characters/goku_normal.webp",
	}

	t.Run("given a valid request, it returns 200", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("SearchCharactersInDatabase", mock.Anything, limit).Return([]domains.Character{characterDomain}, nil)

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.GET("/api/characters/search", SearchCharactersHandler(characterRepoMock))

		req, err := http.NewRequest(http.MethodGet, "/api/characters/search", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("given a valid request, it returns 400 when error limit format", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.GET("/api/characters/search", SearchCharactersHandler(characterRepoMock))

		req, err := http.NewRequest(http.MethodGet, "/api/characters/search?limit=asd", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a valid request, it returns 500 when error unexpected", func(t *testing.T) {
		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("SearchCharactersInDatabase", mock.Anything, limit).Return([]domains.Character{characterDomain}, errors.New("any error"))

		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.GET("/api/characters/search", SearchCharactersHandler(characterRepoMock))

		req, err := http.NewRequest(http.MethodGet, "/api/characters/search", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

}

func Test_DeleteCharacterHandler(t *testing.T) {
	name := "goku"

	t.Run("given a valid request, it returns 200", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()

		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("DeleteCharacterInDatabase", mock.Anything, name).Return(nil)

		r.DELETE("/api/characters/delete/:name", DeleteCharacterHandler(characterRepoMock))

		req, err := http.NewRequest(http.MethodDelete, "/api/characters/delete/goku", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("given a valid request, it returns 500 error domains.ErrCharacterNotDeleted", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()

		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("DeleteCharacterInDatabase", mock.Anything, name).Return(domains.ErrCharacterNotDeleted)

		r.DELETE("/api/characters/delete/:name", DeleteCharacterHandler(characterRepoMock))

		req, err := http.NewRequest(http.MethodDelete, "/api/characters/delete/goku", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("given a valid request, it returns 500 error unexpected", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()

		characterRepoMock := mocks.NewCharacterRepository(t)
		characterRepoMock.On("DeleteCharacterInDatabase", mock.Anything, name).Return(errors.New("any error"))

		r.DELETE("/api/characters/delete/:name", DeleteCharacterHandler(characterRepoMock))

		req, err := http.NewRequest(http.MethodDelete, "/api/characters/delete/goku", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

}
