package domains

import (
	"context"
	"errors"
	"time"
)

var ErrNameIsRequired = errors.New("name is required in json of body")
var ErrCharacterNotFound = errors.New("character not found in external API")
var ErrCharacterNotSave = errors.New("character not save in local database")
var ErrCharacterNotDeleted = errors.New("character not deleted in local database")

type Character struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Ki       int    `json:"ki"`
	Race     string `json:"race"`
	UrlImage string `json:"url_image"`
}

type CharacterRepository interface {
	GetCharacterInExternalAPIByName(
		ctx context.Context,
		name string,
	) (Character, error)
	GetCharacterInDatabaseByName(
		ctx context.Context,
		name string,
	) (Character, error)
	SearchCharactersInDatabase(
		ctx context.Context,
		createdCursor *time.Time,
		limit int,
	) ([]Character, error)
	DeleteCharacterInDatabase(
		ctx context.Context,
		name string,
	) error
}

//go:generate mockery --case=snake --outpkg=mocks --output=./mocks --name=CharacterRepository
