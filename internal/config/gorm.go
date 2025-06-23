package config

import (
	"fmt"
	"seblak-bombom-restful-api/internal/helper"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabaseProd(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("DB_PROD_USERNAME")
	password := viper.GetString("DB_PROD_PASSWORD")
	host := viper.GetString("DB_PROD_HOST")
	port := viper.GetInt("DB_PROD_PORT")
	database_name := viper.GetString("DB_PROD_NAME")
	idleConnection := viper.GetInt("DB_PROD_POOL_IDLE")
	maxConnection := viper.GetInt("DB_PROD_POOL_MAX")
	maxLifeTimeConnection := viper.GetInt("DB_PROD_POOL_LIFETIME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}
	
	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

func NewDatabaseTest(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("DB_TEST_USERNAME")
	password := viper.GetString("DB_TEST_PASSWORD")
	host := viper.GetString("DB_TEST_HOST")
	port := viper.GetInt("DB_TEST_PORT")
	database_name := viper.GetString("DB_TEST_NAME")
	idleConnection := viper.GetInt("DB_TEST_POOL_IDLE")
	maxConnection := viper.GetInt("DB_TEST_POOL_MAX")
	maxLifeTimeConnection := viper.GetInt("DB_TEST_POOL_LIFETIME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}
	
	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

func NewDatabaseDev(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("DB_DEV_USERNAME")
	password := viper.GetString("DB_DEV_PASSWORD")
	host := viper.GetString("DB_DEV_HOST")
	port := viper.GetInt("DB_DEV_PORT")
	database_name := viper.GetString("DB_DEV_NAME")
	idleConnection := viper.GetInt("DB_DEV_POOL_IDLE")
	maxConnection := viper.GetInt("DB_DEV_POOL_MAX")
	maxLifeTimeConnection := viper.GetInt("DB_DEV_POOL_LIFETIME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}
	
	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

func NewDatabaseDocker(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("DB_DOCKER_USERNAME")
	password := viper.GetString("DB_DOCKER_PASSWORD")
	host := viper.GetString("DB_DOCKER_HOST")
	port := viper.GetInt("DB_DOCKER_PORT")
	database_name := viper.GetString("DB_DOCKER_NAME")
	idleConnection := viper.GetInt("DB_DOCKER_POOL_IDLE")
	maxConnection := viper.GetInt("DB_DOCKER_POOL_MAX")
	maxLifeTimeConnection := viper.GetInt("DB_DOCKER_POOL_LIFETIME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}
	
	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

func NewDatabaseDockerTest(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("DB_DOCKER_TEST_USERNAME")
	password := viper.GetString("DB_DOCKER_TEST_PASSWORD")
	host := viper.GetString("DB_DOCKER_TEST_HOST")
	port := viper.GetInt("DB_DOCKER_TEST_PORT")
	database_name := viper.GetString("DB_DOCKER_TEST_NAME")
	idleConnection := viper.GetInt("DB_DOCKER_TEST_POOL_IDLE")
	maxConnection := viper.GetInt("DB_DOCKER_TEST_POOL_MAX")
	maxLifeTimeConnection := viper.GetInt("DB_DOCKER_TEST_POOL_LIFETIME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}
	
	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database : %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
	helper.SaveToLogInfo(args)
}