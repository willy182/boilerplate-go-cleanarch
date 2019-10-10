package config

import (
	"github.com/jinzhu/gorm"

	postgresConfig "github.com/willy182/boilerplate-go-cleanarch/config/postgres"
)

// Config main
type Config struct {
	PostgresDB struct {
		Read, Write *gorm.DB
	}
}

var conf *Config

// Load config
func Load() *Config {
	if conf == nil {
		conf = new(Config)
		conf.PostgresDB.Read = postgresConfig.GetReadDB()
		conf.PostgresDB.Write = postgresConfig.GetWriteDB()
	}

	return conf
}
