package database

import (
	"gorm.io/gorm"
	"os"
)

var (
	DB *gorm.DB
)

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (db *gorm.DB, err error) {
	// Define Database connection variables.
	// Get DB_TYPE value from .env file.
	dbType := os.Getenv("DB_TYPE")
	// Define a new Database connection with right DB type.
	switch dbType {
	case "pgx":
		db, err = PgSqlConnection()
	case "mysql":
		db, err = MysqlConnection()
	}
	if err != nil {

		return nil, err
	}
	DB = db
	return
}
