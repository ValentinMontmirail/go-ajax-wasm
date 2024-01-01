package database

import (
	"context"
	"database/sql"
	embed "wasm/pkg"

	_ "github.com/mattn/go-sqlite3"
)

// MyDB is a struct type that holds a database connection and a context.
type MyDB struct {
	ctx context.Context // ctx is the context, used for managing deadlines, cancellations, etc.
	db  *sql.DB         // db is a pointer to the SQL database connection.
}

// Open opens a database specified by its database driver name and a driver-specific data source name, usually consisting of at least a database name and connection information.
// Most users will open a database via a driver-specific connection helper function that returns a *DB.
// No database drivers are included in the Go standard library.
// See https://golang.org/s/sqldrivers for a list of third-party drivers.
// Open may just validate its arguments without creating a connection to the database.
// The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
// Thus, the Open function should be called just once. It is rarely necessary to close a DB.
func Open(ctx context.Context) (*MyDB, error) {
	// Open a new database connection using the SQLite3 driver, with an in-memory database.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		// If there's an error while opening the database, return nil and the error.
		return nil, err
	}

	// Execute a SQL command (from the 'embed' package) to create necessary tables in the database.
	// The SQL command is executed within the given context 'ctx'.
	if _, err := db.ExecContext(ctx, embed.AuthorSchema); err != nil {
		// If there's an error while executing the SQL command, return nil and the error.
		return nil, err
	}

	// Return a new instance of MyDB with the context and database connection.
	return &MyDB{ctx: ctx, db: db}, nil
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server to finish.
// It is rare to Close a DB, as the DB handle is meant to be long-lived and shared between many goroutines.
func (db *MyDB) Close() error {
	// Call the Close method on the sql.DB object to close the database connection.
	return db.db.Close()
}
