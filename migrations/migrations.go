package migrations

import (
	"fmt"

	"github.com/droplez/droplez-studio/tools/logger"
	postgresClient "github.com/droplez/droplez-studio/third_party/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate() error {
	log := logger.GetServerLogger()
	params := postgresClient.NewConnectionParams()
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", params.Username, params.Password, params.Host, params.Port, params.Database)
	m, err := migrate.New(
		"file://migrations/scripts",
		connectionString)
	if err != nil {
		log.Error(err)
		return err
	}
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info(err)
		} else {
			log.Error(err)
			return err
		}
	}
	return nil
}
