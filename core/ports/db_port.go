package ports

import "database/sql"

// DBPort defines the interface for database operations
type DBPort interface {
	GetConnection() *sql.DB
	Close() error
}
