package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/encilab/dragon-ball/src/domains"
)

type CharacterRepository struct {
	sqlClient     *sql.DB
	clientTimeout time.Duration
}

func NewCharacterRepository(
	sqlClient *sql.DB,
	clientTimeout time.Duration,
) *CharacterRepository {

	return &CharacterRepository{
		sqlClient:     sqlClient,
		clientTimeout: clientTimeout,
	}
}

func (r *CharacterRepository) GetCharacterInExternalAPIByName(
	ctx context.Context,
	name string,
) (domains.Character, error) {
	return domains.Character{}, nil
}

func (r *CharacterRepository) GetCharacterInDatabaseByName(
	ctx context.Context,
	name string,
) (domains.Character, error) {
	return domains.Character{}, nil
}

func (r *CharacterRepository) SearchCharactersInDatabase(
	ctx context.Context,
	createdCursor *time.Time,
	limit int,
) ([]domains.Character, error) {
	return []domains.Character{}, nil
}

func (r *CharacterRepository) DeleteCharacterInDatabase(
	ctx context.Context,
	name string,
) error {
	return nil
}
