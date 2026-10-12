package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/shakerrdev/cs-stalker/internal/db/models"
	"github.com/shakerrdev/cs-stalker/internal/db/models/faceit"
	"github.com/shakerrdev/cs-stalker/internal/db/models/leetify"
	"github.com/shakerrdev/cs-stalker/internal/db/models/steam"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DSNFromEnv() (string, error) {
	keys := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	values := make(map[string]string, len(keys))
	var missing []string
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if !ok {
			missing = append(missing, key)
		}
		values[key] = value
	}
	if len(missing) > 0 {
		return "", fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Europe/Oslo",
		quoteDSN(values["DB_HOST"]), quoteDSN(values["DB_PORT"]), quoteDSN(values["DB_USER"]),
		quoteDSN(values["DB_PASSWORD"]), quoteDSN(values["DB_NAME"])), nil
}

func quoteDSN(value string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(value) + "'"
}

func Connect(dsn string) error {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	if err := Migrate(conn); err != nil {
		return err
	}
	DB = conn
	return nil
}

func Migrate(conn *gorm.DB) error {
	joinTables := []struct {
		model any
		field string
	}{
		{&leetify.LeetifyMatch{}, "Profiles"},
		{&leetify.LeetifyProfile{}, "Matches"},
	}
	for _, jt := range joinTables {
		if err := conn.SetupJoinTable(jt.model, jt.field, &leetify.LeetifyMatchProfile{}); err != nil {
			return fmt.Errorf("set up join table for %T.%s: %w", jt.model, jt.field, err)
		}
	}

	err := conn.AutoMigrate(
		&models.Player{},
		&models.Match{},
		&steam.SteamProfile{},
		&faceit.FaceitProfile{},
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
	)
	if err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}
	return nil
}
