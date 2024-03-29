package storage

import (
	"fmt"
	"log"

	"github.com/satyarth42/chatter/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var userDataDB *gorm.DB

func InitDB(conf config.MySQLConf) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Database)
	userDataDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to establish connection to user data DB err: %s", err)
	}
}

func GetUserDataDB() *gorm.DB {
	return userDataDB
}
