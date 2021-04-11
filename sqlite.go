// Package db provides a kvs for sqlite3.
package db

import (
	"database/sql"
	"fmt"

	"paulfdunn/osh/runtimeh"

	_ "github.com/mattn/go-sqlite3"
)

func Open(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		err = fmt.Errorf("could not open db file, error: %v", err)
		return nil, runtimeh.SourceInfoError("", err)
	}

	return db, nil
}
