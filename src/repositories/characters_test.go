package repositories

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetCharacterInExternalAPIByName(t *testing.T) {
	t.Run("execute get character in external api and success", func(t *testing.T) {
		id := uint(1)
		name := "goku"
		ki := "60.000.000"
		race := "Saiyan"
		image := "https://dragonball-api.com/characters/goku_normal.webp"

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		require.NotNil(t, db)
		require.NotNil(t, mock)

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`INSERT INTO "character_dragonball" ("id", "name", "ki", "race", "image") VALUES ($1, $2, $3, $4, $5)`)).
			WithArgs(id, name, ki, race, image).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := NewCharacterRepository(db, 2000*time.Millisecond)
		characterDomain, err := repo.GetCharacterInExternalAPIByName(context.Background(), name)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)

		assert.Equal(t, name, strings.ToLower(characterDomain.Name))
	})
}

func Test_GetCharacterInDatabaseByName(t *testing.T) {
	t.Run("execute get and success", func(t *testing.T) {
		id := uint(1)
		name := "goku"
		ki := "60.000.000"
		race := "Saiyan"
		image := "https://dragonball-api.com/characters/goku_normal.webp"

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		require.NotNil(t, db)
		require.NotNil(t, mock)

		rows := sqlmock.NewRows([]string{
			"id", "name", "ki", "race", "image",
		}).AddRow(
			id, name, ki, race, image,
		)

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT "id", "name", "ki", "race", "image" FROM "character_dragonball" WHERE "name" = $1`),
		).WithArgs(name).
			WillReturnRows(rows)

		repo := NewCharacterRepository(db, 1000*time.Millisecond)

		auditDomain, err := repo.GetCharacterInDatabaseByName(context.Background(), name)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)

		assert.Equal(t, id, auditDomain.ID)
		assert.Equal(t, name, auditDomain.Name)
		assert.Equal(t, ki, auditDomain.Ki)
		assert.Equal(t, race, auditDomain.Race)
		assert.Equal(t, image, auditDomain.Image)
	})

}

func Test_SearchCharactersInDatabase(t *testing.T) {

	t.Run("execute search and success", func(t *testing.T) {
		id := uint(1)
		name := "goku"
		ki := "60.000.000"
		race := "Saiyan"
		image := "https://dragonball-api.com/characters/goku_normal.webp"

		limit := 10

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		require.NotNil(t, db)
		require.NotNil(t, mock)

		rows := sqlmock.NewRows([]string{"id", "name", "ki", "race", "image"}).
			AddRow(id, name, ki, race, image)

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT id, name, ki, race, image FROM "character_dragonball" ORDER BY id`),
		).WillReturnRows(rows)

		repo := NewCharacterRepository(db, 1000*time.Millisecond)
		results, err := repo.SearchCharactersInDatabase(context.Background(), limit)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

}

func Test_DeleteCharacterInDatabase(t *testing.T) {

	t.Run("execute delete and success", func(t *testing.T) {
		name := "goku"

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		require.NotNil(t, db)
		require.NotNil(t, mock)

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`DELETE FROM "character_dragonball" WHERE "name" = $1`),
		).WithArgs(name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := NewCharacterRepository(db, 1000*time.Millisecond)
		err = repo.DeleteCharacterInDatabase(context.Background(), name)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
	})

}
