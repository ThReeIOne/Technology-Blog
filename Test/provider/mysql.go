package provider

import (
	"Technology-Blog/Test/config"
	"Technology-Blog/Test/log"
	"fmt"
	"time"
)

type Mysql struct {
	DB *sqlx.DB
}

func (db *Mysql) New() *Mysql {
	db.DB = &sqlx.DB{}
	return db
}

func (db *Mysql) Start() {
	db.DB = sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC&time_zone=%%27%%2B00%%3A00%%27&maxAllowedPacket=16777216&readTimeout=60s&writeTimeout=60s&multiStatements=true&charset=utf8mb4,utf8",
		config.Get(config.DBUsername),
		config.Get(config.DBPassword),
		config.Get(config.DBHost),
		config.Get(config.DBPort),
		config.Get(config.DBDatabase),
	))

	if err := db.DB.Ping(); err != nil {
		panic(err)
	}

	db.DB.SetConnMaxLifetime(time.Minute * 9)
	db.DB.SetMaxOpenConns(500)
}

func (db *Mysql) Close() {
	if err := db.DB.Close(); err != nil {
		_ = log.Error(err)
	}
}
