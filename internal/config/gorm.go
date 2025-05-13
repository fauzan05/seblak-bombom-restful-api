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
	username := viper.GetString("database.prod.username")
	password := viper.GetString("database.prod.password")
	host := viper.GetString("database.prod.host")
	port := viper.GetInt("database.prod.port")
	database_name := viper.GetString("database.prod.name")
	idleConnection := viper.GetInt("database.prod.pool.idle")
	maxConnection := viper.GetInt("database.prod.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.prod.pool.lifetime")

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
	username := viper.GetString("database.test.username")
	password := viper.GetString("database.test.password")
	host := viper.GetString("database.test.host")
	port := viper.GetInt("database.test.port")
	database_name := viper.GetString("database.test.name")
	idleConnection := viper.GetInt("database.test.pool.idle")
	maxConnection := viper.GetInt("database.test.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.test.pool.lifetime")

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
	username := viper.GetString("database.dev.username")
	password := viper.GetString("database.dev.password")
	host := viper.GetString("database.dev.host")
	port := viper.GetInt("database.dev.port")
	database_name := viper.GetString("database.dev.name")
	idleConnection := viper.GetInt("database.dev.pool.idle")
	maxConnection := viper.GetInt("database.dev.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.dev.pool.lifetime")

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
	username := viper.GetString("database.docker.username")
	password := viper.GetString("database.docker.password")
	host := viper.GetString("database.docker.host")
	port := viper.GetInt("database.docker.port")
	database_name := viper.GetString("database.docker.name")
	idleConnection := viper.GetInt("database.docker.pool.idle")
	maxConnection := viper.GetInt("database.docker.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.docker.pool.lifetime")

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
	username := viper.GetString("database.docker_test.username")
	password := viper.GetString("database.docker_test.password")
	host := viper.GetString("database.docker_test.host")
	port := viper.GetInt("database.docker_test.port")
	database_name := viper.GetString("database.docker_test.name")
	idleConnection := viper.GetInt("database.docker_test.pool.idle")
	maxConnection := viper.GetInt("database.docker_test.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.docker_test.pool.lifetime")

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