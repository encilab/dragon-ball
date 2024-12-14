package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/encilab/dragon-ball/src/handlers"
	"github.com/encilab/dragon-ball/src/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var folderEnv = "./conf"
var webPort = "8080"

func addRoutes(
	app *gin.Engine,
	sqlClient *sql.DB,
) (*gin.Engine, error) {
	sqlClientTimeout, err := time.ParseDuration(os.Getenv("PSQL_TIMEOUT"))
	if err != nil {
		return app, err
	}

	characterRepository := repositories.NewCharacterRepository(
		sqlClient,
		sqlClientTimeout,
	)

	apiGroup := app.Group("/api")
	apiGroup.GET(
		"/livez",
		handlers.LivezHandler(),
	)

	apiCharacters := apiGroup.Group("/characters")
	apiCharacters.POST(
		"/",
		handlers.GetCharactersHandler(characterRepository),
	)

	return app, nil
}

func newWebApp() (*gin.Engine, error) {
	app := gin.New()

	app.Use(
		gin.Recovery(),
		cors.New(cors.Config{
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowOrigins:     []string{"*"},
			AllowCredentials: true,
		}),
		gin.Logger(),
	)

	return app, nil
}

func validateEnvironmentVariables() error {
	eScope := os.Getenv("SCOPE")
	if eScope != "local" && eScope != "test" && eScope != "qa" && eScope != "prod" {
		return errors.New("there was a problem reading environment variables SCOPE=[local,test,qa,prod]")
	}
	err := godotenv.Load(folderEnv + "/.env." + eScope)
	if err != nil {
		return err
	}

	psqlHost := os.Getenv("PSQL_HOST")
	psqlPort := os.Getenv("PSQL_PORT")
	psqlName := os.Getenv("PSQL_NAME")
	psqlUser := os.Getenv("PSQL_USER")
	psqlPass := os.Getenv("PSQL_PASS")
	psqlTimeout := os.Getenv("PSQL_TIMEOUT")
	if psqlHost == "" || psqlPort == "" || psqlName == "" ||
		psqlUser == "" || psqlPass == "" || psqlTimeout == "" {
		return errors.New("there was a problem reading environment variables PSQL")
	}

	return nil
}

func main() {
	err := validateEnvironmentVariables()
	if err != nil {
		log.Println("error when execute validateEnvironmentVariables, err: " + err.Error())
		return
	}

	app, err := newWebApp()
	if err != nil {
		log.Println("error when execute newWebApp, err: " + err.Error())
		return
	}

	// init sqlClient
	sqlClient, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			os.Getenv("PSQL_HOST"),
			os.Getenv("PSQL_PORT"),
			os.Getenv("PSQL_NAME"),
			os.Getenv("PSQL_USER"),
			os.Getenv("PSQL_PASS"),
			"disable",
		),
	)
	if err != nil {
		log.Println("error when execute sql.Open, err: " + err.Error())
		return
	}
	defer func() {
		err := sqlClient.Close()
		if err != nil {
			log.Println("error when execute sqlClient.Close, err: " + err.Error())
			return
		}
	}()

	app, err = addRoutes(
		app,
		sqlClient,
	)
	if err != nil {
		log.Println("error when execute addRoutes, err: " + err.Error())
		return
	}

	if err := app.Run(
		fmt.Sprintf(":%s", webPort),
	); err != nil {
		log.Println("error when execute app.Run, err: " + err.Error())
		return
	}
}
