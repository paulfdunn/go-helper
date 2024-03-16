package kvs

import (
	"os"
	"path/filepath"
	"testing"
)

func init() {
	t := testing.T{}
	testDir := t.TempDir()
	dataSourceName = filepath.Join(testDir, "test.db")
}

type TestSerialize struct {
	ID   int
	Name string
}

func TestSerializeDeserialize(t *testing.T) {
	testSetupSerialize()

	table := "testTable"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
		return
	}

	tid := 1
	tn := "this is a test..."
	ts := TestSerialize{tid, tn}
	tk := "testKey"
	err = kvs.Serialize(tk, ts)
	if err != nil {
		t.Errorf("Serialize, error: %v", err)
		return
	}
	tsd := TestSerialize{}
	err = kvs.Deserialize(tk, &tsd)
	if err != nil {
		t.Errorf("Deserialize, error: %v", err)
		return
	}
	if tsd.ID != tid || tsd.Name != tn {
		t.Errorf("Deserialize did not work, tsd:%+v", tsd)
		return
	}
}

func TestNegativeDeserialize(t *testing.T) {
	// Deserialize something that was never Serialized.
	testSetupSerialize()

	table := "testTable"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
		return
	}

	tk := "testKey"
	tsd := TestSerialize{ID: 0, Name: ""}
	err = kvs.Deserialize(tk, &tsd)
	if err != nil || tsd.ID != 0 || tsd.Name != "" {
		// expecting: error: scan error: sql: Rows are closed
		t.Errorf("Deserialize, error: %v", err)
		return
	}
}

func testSetupSerialize() error {
	return os.Remove(dataSourceName)
}
