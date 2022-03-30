package database

import (
	"fmt"
	"github.com/channingduan/rpc/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type Database struct {
	databases map[string]*gorm.DB
	config    map[string]config.DatabaseConfig
}

func Register(config *config.Config) *Database {

	database := &Database{
		config:    config.DatabaseConfig,
		databases: make(map[string]*gorm.DB),
	}
	database.initial()

	return database
}

func (db *Database) initial() {
	for name, config := range db.config {
		db.databases[name] = db.connectStart(config)
	}

}

func (db *Database) connectStart(config config.DatabaseConfig) *gorm.DB {

	dbConn := db.connectMaster(config)
	// 主从配置
	sources, replicas := db.connectDialector(config)
	err := dbConn.Use(dbresolver.Register(dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}).SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(100).
		SetMaxOpenConns(200))
	if err != nil {
		fmt.Println("connectStart error: ", err)
	}

	return dbConn
}

func (db *Database) connectMaster(config config.DatabaseConfig) *gorm.DB {

	var dbConn *gorm.DB
	var err error
	var dsn string
	switch config.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)
		dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgresql":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
			config.Host,
			config.Port,
			config.Username,
			config.Password,
			config.Database,
			"Asia/Shanghai",
		)
		dbConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		panic("database driver not found")
	}

	if err != nil {
		panic(fmt.Sprintln("database connect error: ", err))
	}

	return dbConn
}

func (db *Database) connectDialector(config config.DatabaseConfig) ([]gorm.Dialector, []gorm.Dialector) {

	var sources []gorm.Dialector
	var replicas []gorm.Dialector

	if len(config.Sources) > 0 {
		for _, source := range config.Sources {
			switch source.Driver {
			case "MySQL":
				dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
					source.Username,
					source.Password,
					source.Host,
					source.Port,
					source.Database,
				)
				sources = append(sources, mysql.Open(dsn))
			case "PostgreSQL":
				dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
					source.Host,
					source.Port,
					source.Username,
					source.Password,
					source.Database,
					"Asia/Shanghai",
				)
				sources = append(sources, postgres.Open(dsn))
			}
		}
	}

	if len(config.Replicas) > 0 {
		for _, replica := range config.Replicas {
			switch replica.Driver {
			case "MySQL":
				dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
					replica.Username,
					replica.Password,
					replica.Host,
					replica.Port,
					replica.Database,
				)
				replicas = append(replicas, mysql.Open(dsn))
			case "PostgreSQL":
				dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
					replica.Host,
					replica.Port,
					replica.Username,
					replica.Password,
					replica.Database,
					"Asia/Shanghai",
				)
				replicas = append(replicas, postgres.Open(dsn))
			}
		}
	}

	return sources, replicas
}

func (db *Database) AutoMigrate(tables ...interface{}) error {

	return db.NewDatabase().AutoMigrate(tables...)
}

// NewDatabase 单独获取数据库连接
func (db *Database) NewDatabase(name ...string) *gorm.DB {

	conn := "default"
	if len(name) > 0 {
		conn = name[0]
	}

	return db.databases[conn]
}
