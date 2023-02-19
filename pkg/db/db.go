package db

import (
	"errors"
	"learn-im/config"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

func InitMysql() {
	log := logger.New(log.New(os.Stdout, "\n\r", log.LstdFlags), logger.Config{
		SlowThreshold: 500 * time.Millisecond,
		LogLevel:      logger.Info,
		Colorful:      true,
	})
	DB, err := gorm.Open(mysql.Open(config.Conf.Mysql), &gorm.Config{Logger: log})
	if err != nil {
		zap.Error(errors.New("init db error"))
		return
	}
	DB.Debug()
	db, _ := DB.DB()
	// db.SetConnMaxIdleTime()
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisHost,
		Password: config.Conf.RedisPassword,
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		zap.Error(errors.New("redis client error"))
	}
}
