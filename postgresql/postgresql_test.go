package postgresql

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-databasing/connectorsqlite"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-observing/observer"
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

func TestPostgresqlImpl_GetCurrentWatermark(test *testing.T) {
	ctx := context.TODO()
	observer1 := &observer.ObserverNull{
		Id: "Observer 1",
	}
	databaseConnector := &connectorsqlite.Sqlite{
		Filename: "/tmp/sqlite/G2C.db",
	}
	testObject := &PostgresqlImpl{
		LogLevel:          logger.LevelTrace,
		DatabaseConnector: databaseConnector,
	}
	testObject.RegisterObserver(ctx, observer1)
	testObject.GetCurrentWatermark(ctx)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExamplePostgresqlImpl_GetCurrentWatermark() {
	ctx := context.TODO()
	databaseConnector := &connectorsqlite.Sqlite{
		Filename: "/tmp/sqlite/G2C.db",
	}
	testObject := &PostgresqlImpl{
		LogLevel:          logger.LevelTrace,
		DatabaseConnector: databaseConnector,
	}
	testObject.GetCurrentWatermark(ctx)
	// Output:
}