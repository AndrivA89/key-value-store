package transaction

import (
	"database/sql"
	"fmt"

	"github.com/AndrivA89/key-value-store/internal/entity"
)

const tableName = "transactions"

type PostgresParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}

type PostgresLogger struct {
	events chan<- entity.Event
	errors <-chan error
	db     *sql.DB
}

func (l *PostgresLogger) WritePut(key, value string) {
	l.events <- entity.Event{
		EventType: entity.EventPut,
		Key:       key,
		Value:     value,
	}
}

func (l *PostgresLogger) WriteDelete(key string) {
	l.events <- entity.Event{
		EventType: entity.EventDelete,
		Key:       key,
	}
}

func (l *PostgresLogger) Err() <-chan error {
	return l.errors
}

func NewPostgresLogger(config PostgresParams) (Logger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s sslmode=disable password=%s",
		config.Host, config.DbName, config.User, config.Password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresLogger{db: db}

	exists, err := logger.verifyTableExists(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = logger.createTable(tableName); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}

func (l *PostgresLogger) ReadEvents() (<-chan entity.Event, <-chan error) {
	outEvent := make(chan entity.Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		query := `SELECT sequence, event_type, key, value 
					FROM transaction
					ORDER BY sequence`

		rows, err := l.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close()

		e := entity.Event{}
		for rows.Next() {
			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
			}
		}
	}()

	return outEvent, outError
}

func (l *PostgresLogger) Run() {
	events := make(chan entity.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := fmt.Sprintf(`INSERT INTO %s (event_type, key, value) VALUES ($1, $2, $3)`, tableName)
		for e := range events {
			_, err := l.db.Exec(
				query,
				e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *PostgresLogger) verifyTableExists(tableName string) (bool, error) {
	var result string

	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", tableName))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != tableName {
		rows.Scan(&result)
	}

	return result == tableName, rows.Err()
}

func (l *PostgresLogger) createTable(tableName string) error {
	var err error

	createQuery := fmt.Sprintf(`CREATE TABLE %s (
		sequence      BIGSERIAL PRIMARY KEY,
		event_type    SMALLINT,
		key 		  TEXT,
		value         TEXT
	  );`, tableName)

	_, err = l.db.Exec(createQuery)
	if err != nil {
		return err
	}

	return nil
}
