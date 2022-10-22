package main

import (
	"fmt"
	"os"
)

type Env struct {
	Port        string
	DatabaseUrl string
}

const dsn = "postgres://postgres:ppooii12@localhost:5432/shortly"

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
