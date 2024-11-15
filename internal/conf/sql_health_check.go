package conf

import (
	"context"
	"database/sql"
	"fmt"
)

// Wrapper for a SQL connection pool to comply with the health check interface.
type SQLHealthCheck struct {
	DB   *sql.DB
	Name string
}

// PingContext verifies a connection to the database is still alive, establishing a connection if necessary.
func (hc *SQLHealthCheck) PingContext(ctx context.Context) error {
	if hc.DB == nil {
		return fmt.Errorf("failed to ping database, connection is null")
	}

	return hc.DB.PingContext(ctx)
}

// GetName returns the title/name associated with the database connection.
func (hc SQLHealthCheck) GetName() string {
	return hc.Name
}
