package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate(dataSourceName string) error {
	if !strings.Contains(dataSourceName, "://") {
		return fmt.Errorf("store: invalid data source name %q", dataSourceName)
	}

	migrations, err := migrate.New("file://db/migrations", dataSourceName)
	if err != nil {
		return fmt.Errorf("store: failed to initialize migrations: %w", err)
	}

	if err = migrations.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("store: migration failed: %w", err)
	}

	return nil
}
