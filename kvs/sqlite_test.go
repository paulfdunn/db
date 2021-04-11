package kvs

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type kvPair struct {
	key   string
	value []byte
}

var (
	dataSourceName string
)

func init() {
	t := testing.T{}
	testDir := t.TempDir()
	dataSourceName = filepath.Join(testDir, "test.db")
}

func TestDeleteGetSet(t *testing.T) {
	testSetup()

	kvPairs := []kvPair{
		{key: "k1", value: []byte("key1")},
		{key: "k2", value: []byte("key2")},
		{key: "k2", value: []byte("key2.1")},
	}

	kvMap := make(map[string]string)
	table := "testTable"
	kvs, err := New(dataSourceName, table)
	if err != nil {
		t.Errorf("New, error: %v", err)
	}
	for _, v := range kvPairs {
		kvMap[v.key] = string(v.value)
		fmt.Printf("Setting key: %s, value: %s\n", v.key, string(v.value))
		err := kvs.Set(v.key, v.value)
		if err != nil {
			t.Errorf("setting key, error: %v", err)
			return
		}
		value, err := kvs.Get(v.key)
		if err != nil {
			t.Errorf("getting key, error: %v", err)
			return
		}
		if string(value) != string(v.value) {
			t.Errorf("incorrect value")
			return
		}
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

	if count, err := kvs.Delete("k2"); count != 1 || err != nil {
		t.Errorf("deleting key, error: %v", err)
		return
	}
	count, err = rowCount(kvs.dbConn, table)
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

	kvs.DeleteStore()
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

func rowCount(db *sql.DB, table string) (int, error) {
	rows, err := sqlQuery(db, fmt.Sprintf("SELECT * FROM %s;", table))
	if err != nil {
		return 0, fmt.Errorf("getting all rows, error: %v", err)
	}
	count := 0
	for rows.Next() {
		count++
	}
	return count, nil
}

func testSetup() error {
	return os.Remove(dataSourceName)
}
