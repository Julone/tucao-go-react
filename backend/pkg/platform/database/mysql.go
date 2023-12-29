package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tuxiaocao/pkg/logger"
	"tuxiaocao/utils"
)

// MysqlConnection func for connection to Mysql database.
func MysqlConnection() (*gorm.DB, error) {
	// Define database connection settings.
	//maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	//maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	//maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	// Build Mysql connection URL.
	mysqlConnURL, err := utils.ConnectionURLBuilder("mysql")
	if err != nil {
		return nil, err
	}

	// Define database connection for Mysql.
	config := mysql.Config{DSN: mysqlConnURL}
	db, err := gorm.Open(mysql.New(config))
	db = db.Debug()
	logger.Log.Info("fasjdklfjasdklfjaksdl")
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set database connection settings:
	// 	- SetMaxOpenConns: the default is 0 (unlimited)
	// 	- SetMaxIdleConns: defaultMaxIdleConns = 2
	// 	- SetConnMaxLifetime: 0, connections are reused forever
	//db.SetMaxOpenConns(maxConn)
	//db.SetMaxIdleConns(maxIdleConn)
	//db.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database.
	return db, nil
}
