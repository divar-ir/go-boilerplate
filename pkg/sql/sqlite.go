package sql

type SqliteConfig struct {
	FileName string
	InMemory bool
}

func (c SqliteConfig) ConnectionString() string {
	if c.InMemory {
		return ":memory:"
	}

	return c.FileName
}

func (c SqliteConfig) Dialect() string {
	return "sqlite3"
}

func (c SqliteConfig) GetMaxIDLEConnection() int {
	return 1
}

func (c SqliteConfig) GetMaxOpenConnection() int {
	return 1
}
