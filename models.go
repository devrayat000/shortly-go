package main

import "time"

type ShortUrl struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FullUrl   string `json:"fullUrl" gorm:"not null;type:varchar(255);unique;index"`
	ShortUrl  string `json:"shortUrl" gorm:"not null;type:varchar(255);unique;index"`
	ShortKey  string `json:"shortKey" gorm:"not null;type:varchar(255);unique;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
