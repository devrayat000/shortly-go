package main

import (
	"fmt"
	"os"
)

type Env struct {
	Port        string
	DatabaseUrl string
}

const dsn = "host=localhost user=postgres password=ppooii12 dbname=shortly port=5432 sslmode=disable TimeZone=Asia/Shanghai"

func getEnv() *Env {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = dsn
	}

	return &Env{
		Port:        fmt.Sprintf(":%s", port),
		DatabaseUrl: dbUrl,
	}
}
