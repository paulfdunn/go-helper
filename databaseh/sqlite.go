// Package databaseh provides a kvs for sqlite3.
package databaseh

import (
	"database/sql"

	"github.com/paulfdunn/go-helper/osh/runtimeh"

	_ "github.com/mattn/go-sqlite3"
)

func Open(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		err = runtimeh.SourceInfoError("could not open db file", err)
		return nil, runtimeh.SourceInfoError("", err)
	}

	return db, nil
}
