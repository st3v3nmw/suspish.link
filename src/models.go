package main

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Cache *cache.Cache

type FindLinkByURI func(link *Link, URI string) (err error)

type RawLink struct {
	ID        uuid.UUID
	LongURL   string
	SusURI    string
	CreatedAt time.Time
}

type Link struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	LongURL   string    `gorm:"not null;unique;index"`
	SusURI    string    `gorm:"not null;unique;index"`
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

func CachedLink(f FindLinkByURI) FindLinkByURI {
	return func(link *Link, URI string) (err error) {
		raw_link := new(RawLink)
		// check if URI is in the cache
		ctx := context.TODO()
		key := "suspish:" + URI
		if err = Cache.Get(ctx, key, raw_link); err == nil {
			link.ID = raw_link.ID
			link.LongURL = raw_link.LongURL
			link.SusURI = raw_link.SusURI
			link.CreatedAt = raw_link.CreatedAt
			return nil
		}

		// call decorated functin
		if err = f(link, URI); err != nil {
			return err
		}

		// save the result in the cache
		raw_link.ID = link.ID
		raw_link.LongURL = link.LongURL
		raw_link.SusURI = link.SusURI
		raw_link.CreatedAt = link.CreatedAt
		err = Cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: raw_link,
			TTL:   256 * time.Minute,
		})
		if err != nil {
			log.Warn().
				Str("error", err.Error()).
				Msg("Failed to write to the cache")
		}

		return nil
	}
}

func FindLinkBySusURI(link *Link, susURI string) (err error) {
	return DB.Where("sus_uri = ?", susURI).First(link).Error
}

func FindLinkByLongURL(link *Link, longURL string) (err error) {
	return DB.Where("long_url = ?", longURL).First(link).Error
}
