package postgres

import (
	"fmt"
)

func GenerateDataSource(config *PgConfig) string {
	sslMode := "disable"
	if config.SSLMode {
		sslMode = "require"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		sslMode,
	)
}
