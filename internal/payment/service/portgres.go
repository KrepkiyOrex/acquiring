package service

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

type DBConfig struct {
	Driver string `toml:"driver"`
	DSN    string `toml:"dsn"`
}

type Config struct {
	Database DBConfig `toml:"database"`
}

type Application struct {
	db dbInstance
}

type DB struct {
	db *sql.DB
}

type dbInstance interface {
	Close()
	Insert()
}

func (this DB) Close() {
	err := this.db.Close()
	if err != nil {
		log.Fatal("DB close failure: ", err)
	}
}

// postgreSQL DB
func Connect() (*DB, error) {
	// load the configuration from the file
	cfg, err := toml.LoadFile("config/config.toml")
	if err != nil {
		log.Error("Error loading configuration: ", err)
	}

	// getting parameters of connection to the database from the configuration
	dbConfig := cfg.Get("database").(*toml.Tree)
	driver := dbConfig.Get("driver").(string)
	dsn := dbConfig.Get("dsn").(string)

	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		log.Error("Error connecting to the database")
		return &DB{db: sqlDB}, err
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Error("Error pinging the database")
		return &DB{db: sqlDB}, err
	}

	return &DB{db: sqlDB}, nil
}

func GetDB() *DB {
	return &DB{db: &sql.DB{}}
}