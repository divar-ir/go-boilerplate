package sql

type Config interface {
	ConnectionString() string
	Dialect() string
	GetMaxIDLEConnection() int
	GetMaxOpenConnection() int
}

type Migrater interface {
	Migrate() error
}
