package postgres

import (
	"log"

	"github.com/pelletier/go-toml"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Database struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		DBName   string `toml:"dbname"`
		SSLMode  string `toml:"sslmode"`
	} `toml:"database"`
}

func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}
	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = tree.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func Connect() (*gorm.DB, error) {
	config, err := LoadConfig("config/config.toml")
	if err != nil {
		return nil, err
	}

	dsn := "host=" + config.Database.Host +
		" user=" + config.Database.User +
		" password=" + config.Database.Password +
		" dbname=" + config.Database.DBName +
		" port=" + config.Database.Port +
		" sslmode=" + config.Database.SSLMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}
