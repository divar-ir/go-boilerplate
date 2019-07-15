package sql

import "github.com/jinzhu/gorm"

func GetDatabase(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(config.Dialect(), config.ConnectionString())
	if err != nil {
		return nil, err
	}

	maxIDLEConnection := 10
	if config.GetMaxIDLEConnection() != 0 {
		maxIDLEConnection = config.GetMaxIDLEConnection()
	}
	db.DB().SetMaxIdleConns(maxIDLEConnection)

	maxOpenConnection := 100
	if config.GetMaxOpenConnection() != 0 {
		maxOpenConnection = config.GetMaxOpenConnection()
	}

	db.DB().SetMaxIdleConns(maxOpenConnection)

	return db, nil
}
