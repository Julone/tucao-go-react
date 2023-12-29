package setup

import (
	"fmt"
	"time"
	"tuxiaocao/pkg/logger"
	"tuxiaocao/pkg/platform/database"
	models2 "tuxiaocao/routes/models"
)

func ReadyDatabase() {
	// database init
	db, err := database.OpenDBConnection()
	if err != nil {
		logger.Log.Panicf("mysql is error %v", err)
	}
	err = db.AutoMigrate(models2.Product{}, models2.User{})
	if err != nil {
		logger.Log.Errorf("mysql migrate is error %v", err)
	}

}

func ReadyLogger() {
	logger.Replace(logger.New("./", "rmp", "debug", time.Hour*240, time.Hour*24, 100))
	if logger.Log == nil {
		panic(fmt.Errorf("log start error"))
	}
}

func InitAll(Components ...string) {
	if len(Components) == 0 {
		Components = []string{"logger", "mysql"}
	}
	for _, val := range Components {
		switch val {
		case "logger":
			ReadyLogger()
		case "mysql":
			ReadyDatabase()
		}
	}
}
