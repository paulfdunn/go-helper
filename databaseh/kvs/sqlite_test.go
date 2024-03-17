package kvs

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paulfdunn/go-helper/osh/runtimeh"
)

type kvPairs []kvPair
type kvPair struct {
	key   string
	value []byte
}

var (
	dataSourceName string
	kvps           kvPairs
)

func init() {
	t := testing.T{}
	testDir := t.TempDir()
	dataSourceName = filepath.Join(testDir, "test.db")
}

func TestDeleteGetSet(t *testing.T) {
	testSetup()

	table := "testTable"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
	}

	kvMap, err := kvps.add(t, kvs)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	// Close and re-open to show that New will open an existing database.
	kvs.Close()
	kvs, err = New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
	}

	if count, err := kvs.Delete("k2"); count != 1 || err != nil {
		t.Errorf("deleting key, error: %v", err)
		return
	}
	count, err := rowCount(kvs.dbConn, table)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if count != len(kvMap)-1 {
		t.Errorf(fmt.Sprintf("wrong number of rows:%d", count))
		return
	}

	if count, err := kvs.Delete("k1"); count != 1 || err != nil {
		t.Errorf("deleting key, error: %v", err)
		return
	}
	count, err = rowCount(kvs.dbConn, table)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if count != len(kvMap)-2 {
		t.Errorf(fmt.Sprintf("wrong number of rows:%d", count))
		return
	}
}

func TestKeys(t *testing.T) {
	testSetup()

	table := "testTable"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
	}

	kvMap, err := kvps.add(t, kvs)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	keys, err := kvs.Keys()
	if err != nil || len(keys) != len(kvMap) {
		t.Errorf(fmt.Sprintf("error or wrong number of keys:%v", keys))
		return
	}
	for _, k := range keys {
		if _, ok := kvMap[k]; !ok {
			t.Errorf(fmt.Sprintf("key not in kvMap:%s", k))
			return
		}
	}
}
func TestRowCount(t *testing.T) {
	testSetup()

	table := "testTable"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
	}

	kvMap, err := kvps.add(t, kvs)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	count, err := rowCount(kvs.dbConn, table)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if count != len(kvMap) {
		t.Errorf(fmt.Sprintf("wrong number of rows:%d", count))
		return
	}
}

// Deleting a non-existent key does not produce an error, but the count is zero.
func TestDeleteNegative(t *testing.T) {
	testSetup()

	table := "testTable"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
		return
	}

	count, err := kvs.Delete("k1")
	if err != nil || count != 0 {
		t.Error("Delete with no data did not produce error or had non-zero count")
		return
	}
}

func TestDeleteStore(t *testing.T) {
	testSetup()

	table := "testTableN2"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
		return
	}

	err = kvs.DeleteStore()
	if err != nil {
		t.Errorf("New, error: %v", err)
		return
	}
	// The key doesn't matter as the table was deleted.
	_, err = kvs.Get("")
	if err == nil {
		// Should produce: "error: no such table: testTableN2"
		t.Errorf("no error Getting when table was deleted.")
		return
	}
}

func TestGetNegative(t *testing.T) {
	testSetup()

	table := "testTableN1"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
	}

	b, err := kvs.Get("k1")
	if !(b == nil && err == nil) {
		t.Error("Get with invalid key should produce no data and no error")
		return
	}
}

func (kvps kvPairs) add(t *testing.T, kvs KVS) (map[string]string, error) {
	kvMap := make(map[string]string)
	for _, v := range kvps {
		kvMap[v.key] = string(v.value)
		fmt.Printf("kvPairs.add setting key: %s, value: %s\n", v.key, string(v.value))
		err := kvs.Set(v.key, v.value)
		if err != nil {
			t.Errorf("setting key, error: %v", err)
			return nil, err

		}
		value, err := kvs.Get(v.key)
		if err != nil {
			t.Errorf("getting key, error: %v", err)
			return nil, err
		}
		if string(value) != string(v.value) {
			t.Errorf("incorrect value")
			return nil, err
		}
	}
	return kvMap, nil
}

func rowCount(db *sql.DB, table string) (int, error) {
	rows, err := sqlQuery(db, fmt.Sprintf("SELECT * FROM %s;", table))
	if err != nil {
		return 0, runtimeh.SourceInfoError("getting all rows", err)
	}
	count := 0
	for rows.Next() {
		count++
	}
	return count, nil
}

func testSetup() error {
	kvps = kvPairs{
		{key: "k1", value: []byte("key1")},
		{key: "k2", value: []byte("key2")},
		{key: "k2", value: []byte("key2.1")},
		{key: "k3", value: []byte("key3")},
		{key: "k4", value: []byte("key4")},
	}
	return os.Remove(dataSourceName)
}
