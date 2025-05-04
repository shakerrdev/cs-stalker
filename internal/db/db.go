package db

import (
	"fmt"
	"log"
	"os"

	"github.com/finnbechinka/cs-stalker/internal/db/models"
	"github.com/finnbechinka/cs-stalker/internal/db/models/faceit"
	"github.com/finnbechinka/cs-stalker/internal/db/models/leetify"
	"github.com/finnbechinka/cs-stalker/internal/db/models/steam"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Retrieve necessary environment variables
	usr, exists := os.LookupEnv("DB_USER")
	if !exists {
		log.Fatalf("missing environment variable: DB_USER")
	}
	pw, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		log.Fatalf("missing environment variable: DB_PASSWORD")
	}
	host, exists := os.LookupEnv("DB_HOST")
	if !exists {
		log.Fatalf("missing environment variable: DB_HOST")
	}
	port, exists := os.LookupEnv("DB_PORT")
	if !exists {
		log.Fatalf("missing environment variable: DB_PORT")
	}
	dbname, exists := os.LookupEnv("DB_NAME")
	if !exists {
		log.Fatalf("missing environment variable: DB_NAME")
	}

	// Construct the database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Oslo",
		host, usr, pw, dbname, port)

	// Open a connection to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	} else {
		log.Println("connected to database")
	}

	// Run AutoMigrate to update the schema
	models := []interface{}{
		&models.Player{},
		&models.Match{},

		// Leetify models
		&leetify.LeetifyProfile{},
		&leetify.LeetifyClub{},
		&leetify.LeetifyTeammate{},
		&leetify.LeetifyRank{},
		&leetify.LeetifyHighlight{},
		&leetify.LeetifyPersonalBest{},
		&leetify.LeetifyMatch{},
		&leetify.LeetifyMatchDetails{},
		&leetify.LeetifyPlayerStats{},
		&leetify.LeetifyAgent{},
		&leetify.LeetifyParty{},

		// Other services
		&faceit.FaceitProfile{},
		&steam.SteamProfile{},
	}

	migration_error := false

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			log.Printf("failed to migrate %T: %v", model, err)
			migration_error = true
		}
	}
	if migration_error {
		log.Fatalf("failed to auto migrate")
	} else {
		log.Println("auto migrated database")
	}
}
