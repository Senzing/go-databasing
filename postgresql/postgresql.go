package postgresql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// PostgresqlImpl is the default implementation of the SqlExecutor interface.
type PostgresqlImpl struct {
	DatabaseConnector driver.Connector
	isTrace           bool
	logger            messagelogger.MessageLoggerInterface
	LogLevel          logger.Level
	observers         subject.Subject
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (sqlexecutor *PostgresqlImpl) getLogger() messagelogger.MessageLoggerInterface {
	if sqlexecutor.logger == nil {
		sqlexecutor.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return sqlexecutor.logger
}

// Notify registered observers.
func (sqlexecutor *PostgresqlImpl) notify(ctx context.Context, messageId int, err error, details map[string]string) {
	now := time.Now()
	details["subjectId"] = strconv.Itoa(ProductId)
	details["messageId"] = strconv.Itoa(messageId)
	details["messageTime"] = strconv.FormatInt(now.UnixNano(), 10)
	if err != nil {
		details["error"] = err.Error()
	}
	message, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		sqlexecutor.observers.NotifyObservers(ctx, string(message))
	}
}

// Trace method entry.
func (sqlexecutor *PostgresqlImpl) traceEntry(errorNumber int, details ...interface{}) {
	sqlexecutor.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (sqlexecutor *PostgresqlImpl) traceExit(errorNumber int, details ...interface{}) {
	sqlexecutor.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The GetCurrentWatermark does a database call for each line scanned.

Input
  - ctx: A context to control lifecycle.
*/
func (sqlexecutor *PostgresqlImpl) GetCurrentWatermark(ctx context.Context) (string, int, error) {
	var (
		oid  string
		age  int
		size string
	)

	// Entry tasks.

	if sqlexecutor.isTrace {
		sqlexecutor.traceEntry(1)
	}
	entryTime := time.Now()
	sqlStatement := "SELECT c.oid::regclass, age(c.relfrozenxid), pg_size_pretty(pg_total_relation_size(c.oid)) FROM pg_class c JOIN pg_namespace n on c.relnamespace = n.oid WHERE relkind IN ('r', 't', 'm') AND n.nspname NOT IN ('pg_toast') ORDER BY 2 DESC LIMIT 1;"

	// Open a database connection.

	database := sql.OpenDB(sqlexecutor.DatabaseConnector)
	defer database.Close()
	err := database.PingContext(ctx)
	if err != nil {
		return "", 0, err
	}

	// Get the Row.

	row := database.QueryRowContext(ctx, sqlStatement)
	err = row.Scan(&oid, &age, &size)
	if err != nil {
		return "", 0, err
	}

	// Exit tasks.

	if sqlexecutor.observers != nil {
		go func() {
			details := map[string]string{
				"oid": oid,
				"age": strconv.Itoa(age),
			}
			sqlexecutor.notify(ctx, 8001, err, details)
		}()
	}
	if sqlexecutor.isTrace {
		defer sqlexecutor.traceExit(2, oid, age, err, time.Since(entryTime))
	}
	return oid, age, err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (sqlexecutor *PostgresqlImpl) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if sqlexecutor.isTrace {
		sqlexecutor.traceEntry(3, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if sqlexecutor.observers == nil {
		sqlexecutor.observers = &subject.SubjectImpl{}
	}
	err := sqlexecutor.observers.RegisterObserver(ctx, observer)
	if sqlexecutor.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverId(ctx),
			}
			sqlexecutor.notify(ctx, 8002, err, details)
		}()
	}
	if sqlexecutor.isTrace {
		defer sqlexecutor.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (sqlexecutor *PostgresqlImpl) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if sqlexecutor.isTrace {
		sqlexecutor.traceEntry(5, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	sqlexecutor.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	sqlexecutor.isTrace = (sqlexecutor.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if sqlexecutor.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			sqlexecutor.notify(ctx, 8003, err, details)
		}()
	}
	if sqlexecutor.isTrace {
		defer sqlexecutor.traceExit(6, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (sqlexecutor *PostgresqlImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if sqlexecutor.isTrace {
		sqlexecutor.traceEntry(7, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if sqlexecutor.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		sqlexecutor.notify(ctx, 8004, err, details)
	}
	err = sqlexecutor.observers.UnregisterObserver(ctx, observer)
	if !sqlexecutor.observers.HasObservers(ctx) {
		sqlexecutor.observers = nil
	}
	if sqlexecutor.isTrace {
		defer sqlexecutor.traceExit(8, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
