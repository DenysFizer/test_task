package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"test_task/config"
)

var (
	DB *pgxpool.Pool
)

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	var err error
	DB, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	err = DB.Ping(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Unable to ping to database: %v", err))
	}
	log.Println("Db Connected")
}

func CloseDB() {
	DB.Close()
}
