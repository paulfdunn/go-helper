// Package kvs implements a key/value store. This implementation is for SQLITE3.
package kvs

import (
	"database/sql"
	"fmt"

	"github.com/paulfdunn/go-helper/databaseh"
	"github.com/paulfdunn/go-helper/osh/runtimeh"
)

// KVS is an instance for key/value storage.
type KVS struct {
	dbConn *sql.DB
	table  string
}

// New creates a new key/value store, with a new or existing table, in the database for key/value storage.
// The database file is created if it does not exist; an existing file is used if present.
// The returned object has a single connection. If performance is an issue, create a pool of connections.
// The GO sql package insures single threaded access to the connection, and thus it is thread safe.
func New(dbConnectionString string, table string) (KVS, error) {
	var dbConn *sql.DB
	var err error
	if dbConn, err = databaseh.Open(dbConnectionString); err != nil {
		return KVS{}, runtimeh.SourceInfoError("opening db", err)
	}

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (key string NOT NULL PRIMARY KEY, value BLOB);`, table)
	_, err = sqlExec(dbConn, query)
	return KVS{dbConn, table}, runtimeh.SourceInfoError("creating kvs table", err)
}

// Close closes the database connection.
func (kvs KVS) Close() error {
	return kvs.dbConn.Close()
}

// Delete deletes a key from the KVS; returns the count, which is zero (and no error) if the key did not exist.
func (kvs KVS) Delete(key string) (int64, error) {
	if kvs.dbConn == nil {
		return 0, fmt.Errorf("%s kvs dbConn is nil", runtimeh.SourceInfo())
	}

	stmt, err := kvs.dbConn.Prepare(fmt.Sprintf(`DELETE FROM %s WHERE key=?;`, kvs.table))
	if err != nil {
		return 0, runtimeh.SourceInfoError("", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(key)
	if err != nil {
		return 0, runtimeh.SourceInfoError("", err)
	}

	var count int64
	if count, err = res.RowsAffected(); err != nil {
		return 0, err
	}

	return count, nil
}

// DeleteStore drops the table associated with the KVS.
func (kvs KVS) DeleteStore() error {
	if kvs.dbConn == nil {
		return fmt.Errorf("%s kvs db is nil", runtimeh.SourceInfo())
	}

	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s ;`, kvs.table)
	_, err := sqlExec(kvs.dbConn, query)
	return runtimeh.SourceInfoError("", err)
}

// Get gets a value from the KVS.
// If the return data and error are both nil, the key did not exist.
func (kvs KVS) Get(key string) ([]byte, error) {
	if kvs.dbConn == nil {
		return nil, fmt.Errorf("%s kvs db is nil", runtimeh.SourceInfo())
	}

	stmt, err := kvs.dbConn.Prepare(fmt.Sprintf(`SELECT value FROM %s WHERE key=?;`, kvs.table))
	if err != nil {
		return nil, runtimeh.SourceInfoError("", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(key)
	if err != nil {
		return nil, runtimeh.SourceInfoError("", err)
	}
	defer rows.Close()

	rslt := rows.Next()
	if !rslt && rows.Err() == nil {
		// The key did not exist
		return nil, runtimeh.SourceInfoError("", err)
	}
	var value []byte
	err = rows.Scan(&value)
	if err != nil {
		return nil, runtimeh.SourceInfoError("scan error", err)
	}

	err = rows.Err()
	if err != nil {
		return nil, runtimeh.SourceInfoError("scan iteration error", err)
	}

	return value, nil
}

// Keys returns all keys in the store.
func (kvs KVS) Keys() ([]string, error) {
	if kvs.dbConn == nil {
		return nil, fmt.Errorf("%s kvs db is nil", runtimeh.SourceInfo())
	}

	rows, err := sqlQuery(kvs.dbConn, fmt.Sprintf("SELECT key FROM %s;", kvs.table))
	if err != nil {
		return nil, runtimeh.SourceInfoError("getting all keys", err)
	}
	keys := []string{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, runtimeh.SourceInfoError("scanning all keys", err)
		}
		keys = append(keys, key)
	}
	return keys, nil
}

// Set sets a value for the specified key in the KVS.
func (kvs KVS) Set(key string, value []byte) error {
	if kvs.dbConn == nil {
		return fmt.Errorf("%s kvs db is nil", runtimeh.SourceInfo())
	}

	stmt, err := kvs.dbConn.Prepare(fmt.Sprintf(`REPLACE INTO %s Values (?,?);`, kvs.table))
	if err != nil {
		return runtimeh.SourceInfoError("", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(key, value)
	if err != nil {
		return runtimeh.SourceInfoError("", err)
	}

	return nil
}

func sqlExec(db *sql.DB, query string) (sql.Result, error) {
	result, err := db.Exec(query)
	// if err != nil {
	// 	fmt.Printf("sqlExec error: %v, %s\n", err, query)
	// }
	// lastID, _ := result.LastInsertId()
	// rowsAffected, _ := result.RowsAffected()
	// fmt.Printf("lastID: %d, rowsAffected: %d\n", lastID, rowsAffected)
	return result, runtimeh.SourceInfoError("", err)
}

func sqlQuery(db *sql.DB, query string) (*sql.Rows, error) {
	rows, err := db.Query(query)
	// if err != nil {
	// 	fmt.Printf("sqlQuery error: %v, %s\n", err, query)
	// }
	return rows, runtimeh.SourceInfoError("", err)
}
