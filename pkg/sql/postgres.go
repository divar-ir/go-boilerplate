package sql

import "fmt"

type PostgresConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	Database          string
	SSL               bool
	MaxIDLEConnection int
	MaxOpenConnection int
}

func (c PostgresConfig) ConnectionString() string {
	sslMode := "enable"
	if !c.SSL {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s binary_parameters=yes",
		c.Host,
		c.Port,
		c.Username,
		c.Database,
		c.Password,
		sslMode,
	)
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) GetMaxIDLEConnection() int {
	return c.MaxIDLEConnection
}

func (c PostgresConfig) GetMaxOpenConnection() int {
	return c.MaxOpenConnection
}
