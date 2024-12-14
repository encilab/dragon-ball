package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetCharacterInExternalAPIByName(t *testing.T) {
	t.Run("execute get character in external api and success", func(t *testing.T) {
		name := "goku"

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		require.NotNil(t, db)
		require.NotNil(t, mock)

		repo := NewCharacterRepository(db, 1000*time.Millisecond)
		characterDomain, err := repo.GetCharacterInExternalAPIByName(context.Background(), name)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)

		assert.Equal(t, name, characterDomain.Name)
	})
}
