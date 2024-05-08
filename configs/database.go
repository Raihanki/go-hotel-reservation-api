package configs

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDatabase() {
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		ENV.DB_USERNAME,
		ENV.DB_PASSWORD,
		ENV.DB_HOST,
		ENV.DB_PORT,
		ENV.DB_DATABASE,
	)

	db, errDatabase := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errDatabase != nil {
		log.Fatal("Error connecting to database")
	}

	DB = db
}
