package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Link struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	LongURL   string    `gorm:"not null,unique"`
	SusURI    string    `gorm:"not null,unique"`
	Clicks    uint      `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func CreateLink(link *Link) (err error) {
	return DB.Create(link).Error
}

func ListLinks(links *[]Link) (err error) {
	return DB.Find(links).Error
}

func FindLinkBySusURI(link *Link, susURI string) (err error) {
	return DB.Where("sus_uri = ?", susURI).First(link).Error
}

func FindLinkByLongURL(link *Link, longURL string) (err error) {
	return DB.Where("long_url = ?", longURL).First(link).Error
}

func UpdateLink(link *Link) (err error) {
	return DB.Save(link).Error
}
