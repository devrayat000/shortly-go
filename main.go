package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=postgres password=ppooii12 dbname=shortly port=5432 sslmode=disable TimeZone=Asia/Shanghai"

type ShortUrl struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FullUrl   string `json:"fullUrl" gorm:"not null;type:varchar(255);unique;index"`
	ShortUrl  string `json:"shortUrl" gorm:"not null;type:varchar(255);unique;index"`
	ShortKey  string `json:"shortKey" gorm:"not null;type:varchar(255);unique;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/shorten", func(ctx *gin.Context) {
		url := ctx.PostForm("url")
		uuid := uniuri.NewLen(7)
		trimmedUrl := fmt.Sprintf("http://localhost:8080/%s", uuid)
		shortUrl := ShortUrl{
			FullUrl:  url,
			ShortUrl: trimmedUrl,
			ShortKey: uuid,
		}

		result := db.Model(&ShortUrl{}).Where(&ShortUrl{FullUrl: url}, "FullUrl").FirstOrCreate(&shortUrl)

		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, shortUrl)
	})

	r.GET("/:shortUrl", func(ctx *gin.Context) {
		code := ctx.Param("shortUrl")
		var shortUrl ShortUrl
		result := db.Model(&ShortUrl{}).Where(&ShortUrl{ShortKey: code}).Select("FullUrl").First(&shortUrl)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No url found for this code!"})
		}
		fmt.Print(shortUrl)
		ctx.Redirect(http.StatusPermanentRedirect, shortUrl.FullUrl)
	})

	return r
}

func main() {
	ConfigRuntime()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database.")
	}
	db.AutoMigrate(&ShortUrl{})

	r := setupRouter(db)
	// Listen and Server in 0.0.0.0:8080
	err = r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
