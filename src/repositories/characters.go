package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	client := &http.Client{
		Timeout: r.clientTimeout,
	}
	url := fmt.Sprintf("https://dragonball-api.com/api/characters?name=%s", name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domains.Character{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return domains.Character{}, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domains.Character{}, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var characters []domains.Character
	if err := json.NewDecoder(resp.Body).Decode(&characters); err != nil {
		return domains.Character{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(characters) == 0 {
		return domains.Character{}, domains.ErrCharacterNotFoundInExternalAPI
	}

	characters[0].Name = strings.ToLower(characters[0].Name)
	err = r.setCharacterInDatabase(ctx, characters[0].ID, characters[0].Name, characters[0].Ki, characters[0].Race, characters[0].Image)
	if err != nil {
		return domains.Character{}, err
	}

	return characters[0], nil
}

func (r *CharacterRepository) setCharacterInDatabase(
	ctx context.Context,
	id uint,
	name,
	ki,
	race,
	image string,
) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.clientTimeout)
	defer cancel()

	tx, err := r.sqlClient.BeginTx(ctxTimeout, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println(rbErr)
			}
		}
	}()

	var query string
	var args []interface{}

	query = `INSERT INTO "character_dragonball" ("id", "name", "ki", "race", "image") VALUES ($1, $2, $3, $4, $5)`
	args = []interface{}{
		id,
		strings.ToLower(name),
		ki,
		race,
		image,
	}

	_, err = tx.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			err = domains.ErrCharacterAlreadyExistInDatabase
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CharacterRepository) GetCharacterInDatabaseByName(
	ctx context.Context,
	name string,
) (domains.Character, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.clientTimeout)
	defer cancel()

	var character domains.Character
	err := r.sqlClient.QueryRowContext(
		ctxTimeout,
		`SELECT "id", "name", "ki", "race", "image" FROM "character_dragonball" WHERE "name" = $1`,
		name,
	).Scan(
		&character.ID,
		&character.Name,
		&character.Ki,
		&character.Race,
		&character.Image,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domains.Character{}, domains.ErrCharacterNotFoundInDatabase
		}
		return domains.Character{}, err
	}

	return character, nil
}

func (r *CharacterRepository) SearchCharactersInDatabase(
	ctx context.Context,
	limit int,
) ([]domains.Character, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.clientTimeout)
	defer cancel()

	query := `SELECT id, name, ki, race, image FROM "character_dragonball"`
	var args []interface{}
	query += " ORDER BY id"

	rows, err := r.sqlClient.QueryContext(ctxTimeout, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []domains.Character{}
	for rows.Next() {
		if len(results) >= limit {
			break
		}

		var character domains.Character
		if err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.Ki,
			&character.Race,
			&character.Image,
		); err != nil {
			return nil, err
		}

		results = append(results, character)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *CharacterRepository) DeleteCharacterInDatabase(
	ctx context.Context,
	name string,
) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.clientTimeout)
	defer cancel()

	tx, err := r.sqlClient.BeginTx(ctxTimeout, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Println(err)
			}
		}
	}()

	result, err := tx.ExecContext(
		ctxTimeout,
		`DELETE FROM "character_dragonball" WHERE "name" = $1`,
		strings.ToLower(name),
	)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return domains.ErrCharacterNotDeleted
	}

	return nil
}
