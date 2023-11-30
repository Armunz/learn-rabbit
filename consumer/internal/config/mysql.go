package config

import (
	"database/sql"
	"log"
	"time"
)

func NewMySQLDatabase(config ConsumerConfig) *sql.DB {
	db, err := sql.Open("mysql", config.DBUserDSN)
	if err != nil {
		log.Fatalln("failed to connect to mysql database, ", err)
	}

	db.SetMaxOpenConns(config.MYSQLMaxConn)
	db.SetConnMaxLifetime(time.Duration(config.ConnLifeTimeSecond) * time.Second)
	db.SetMaxIdleConns(config.MYSQLIdleConn)

	if err := db.Ping(); err != nil {
		log.Fatalln("failed to ping mysql, ", err)
	}

	return db
}

func CleanUp(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalln("failed to close mysql connection, ", err)
	}
}
