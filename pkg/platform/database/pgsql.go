package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"tuxiaocao/pkg/logger"
	"tuxiaocao/utils"
)

// MysqlConnection func for connection to Mysql database.
func PgSqlConnection() (*gorm.DB, error) {
	// Define database connection settings.
	//maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	//maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	//maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	// Build Mysql connection URL.
	mysqlConnURL, err := utils.ConnectionURLBuilder("postgres")
	logger.Log.Info("connect", mysqlConnURL)
	if err != nil {
		return nil, err
	}

	config := postgres.Config{DSN: mysqlConnURL}
	db, err := gorm.Open(postgres.New(config))

	// Define database connection for Mysql.
	//config := mysql.Config{DSN: mysqlConnURL}
	//db, err := gorm.Open(mysql.New(config))
	db = db.Debug()
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}
	return db, nil
}
