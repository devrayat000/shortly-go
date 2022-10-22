package main

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func StartWorkers() {
	go statsWorker()
}

func StartDatabase() *gorm.DB {
	env := getEnv()
	db, err := gorm.Open(postgres.Open(env.DatabaseUrl), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Could not connect to database - Error: %v", err))
	}
	db.AutoMigrate(&ShortUrl{})
	return db
}

func StartServer(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/shorten", CreateShortUrl(db))
	r.GET("/:shortUrl", RetrieveShortUrl(db))

	env := getEnv()
	if err := r.Run(env.Port); err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}

	return r
}

func main() {
	ConfigRuntime()
	StartWorkers()
	db := StartDatabase()
	StartServer(db)
}
