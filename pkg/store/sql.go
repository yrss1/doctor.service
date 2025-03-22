package store

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SQLX struct {
	Client *sqlx.DB
}

func New(dbSource string) (store SQLX, err error) {
	driverName := strings.ToLower(strings.Split(dbSource, "://")[0])
	store.Client, err = sqlx.Connect(driverName, dbSource)
	if err != nil {
		return store, fmt.Errorf("failed to connect to database: %w", err)
	}
	store.Client.SetMaxOpenConns(20)
	store.Client.SetMaxIdleConns(10)
	return
}
