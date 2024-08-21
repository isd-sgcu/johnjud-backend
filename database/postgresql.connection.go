package database

import (
	"github.com/isd-sgcu/johnjud-gateway/config"
	"github.com/isd-sgcu/johnjud-gateway/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitPostgresDatabase(conf *config.Database, isDebug bool) (db *gorm.DB, err error) {
	gormConf := &gorm.Config{TranslateError: true}

	if !isDebug {
		gormConf.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err = gorm.Open(postgres.Open(conf.Url), gormConf)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.User{}, &model.AuthSession{})
	if err != nil {
		return nil, err
	}

	return
}
