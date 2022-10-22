package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateShortUrl(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
	}
}

func RetrieveShortUrl(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Param("shortUrl")
		var shortUrl ShortUrl
		result := db.Model(&ShortUrl{}).Where(&ShortUrl{ShortKey: code}).Select("FullUrl").First(&shortUrl)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "No url found for this code!"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			}
			return
		}

		ctx.Redirect(http.StatusPermanentRedirect, shortUrl.FullUrl)
	}
}
