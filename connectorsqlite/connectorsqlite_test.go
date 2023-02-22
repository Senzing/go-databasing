package connectorsqlite

import (
	"context"
	"fmt"
	"os"
	"testing"
)

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error = nil
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSqlite_Connect(test *testing.T) {
	ctx := context.TODO()
	databaseConnector, err := NewConnector(ctx, "/tmp/sqlite/G2C.db")
	if err != nil {
		test.Fatal(err)

	}
	databaseConnector.Connect(ctx)
}
