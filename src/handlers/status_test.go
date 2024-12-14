package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LivezHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/livez", LivezHandler())

	t.Run("given a valid request, it returns 200", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/api/livez", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

}

func Test_ReadyzHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	require.NotNil(t, db)
	require.NotNil(t, mock)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/readyz", ReadyzHandler(db))

	t.Run("given a valid request, it returns 200", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/api/readyz", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

}
