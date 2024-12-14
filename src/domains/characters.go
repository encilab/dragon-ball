package domains

import (
	"context"
	"errors"
)

var ErrNameIsRequired = errors.New("name is required in json of body")
var ErrCharacterNotFoundInExternalAPI = errors.New("character not found in external API")
var ErrCharacterNotFoundInDatabase = errors.New("character not found in database")
var ErrCharacterAlreadyExistInDatabase = errors.New("character already exist in database")
var ErrCharacterNotSave = errors.New("character not save in local database")
var ErrCharacterNotDeleted = errors.New("character not deleted in local database")

type Character struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Ki    string `json:"ki"`
	Race  string `json:"race"`
	Image string `json:"image"`
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
		limit int,
	) ([]Character, error)
	DeleteCharacterInDatabase(
		ctx context.Context,
		name string,
	) error
}

//go:generate mockery --case=snake --outpkg=mocks --output=./mocks --name=CharacterRepository
