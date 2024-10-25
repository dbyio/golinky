package urlshortener

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

const prefix = "urlshortener"

type Shortenerer interface {
	CreateShortUrl(string) (string, error)
	UseShortUrl(string) (string, error)
	HasCallbackUrl() bool
}

var _ Shortenerer = Shortener{}

type Shortener struct {
	Length      uint8  // number of characters in the short URL's path
	Timeout     int    // lifetime of the short URL (in seconds)
	RedisUrl    string // redis server url (rediss://user:pass@host:port)
	BaseUrl     string // base url of the shortener
	Seed        string // secret see for the hash function
	CallbackUrl string // callback endpoint (no default value)
	storage     *redisStore
	ctx         context.Context
}

type redisStore struct {
	rdb *redis.Client
}

type Stats struct {
	DbSize int64 `json:"dbsize"`
}

func New(file string) *Shortener {
	s := loadConfig(file)
	if s.BaseUrl == "" {
		log.Fatal("Failed to load configuration, please verify file format.")
	}

	// minimum value for length and timeout
	if s.Length < 5 {
		s.Length = 5
	}
	if s.Timeout < 3 {
		s.Timeout = 3
	}

	(*s).ctx = context.Background()
	(*s).storage = initStore(s.RedisUrl, s.ctx)

	return s
}

func (s Shortener) CreateShortUrl(longUrl string) (string, error) {
	shortUrl := s.BaseUrl + shrinkUrl(longUrl, s.Seed, s.Length)

	err := s.store(shortUrl, longUrl)
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

func (s Shortener) UseShortUrl(path string) (string, error) {
	shortUrl := s.BaseUrl + path
	longUrl, err := s.fetch(shortUrl)
	if err != nil {
		return "", err
	}
	return longUrl, nil
}

func shrinkUrl(url string, seed string, length uint8) string {
	hash := hashMe(url, seed)
	l := min(len(hash), int(length))
	return hash[:l]
}

// Storage
func initStore(redisUrl string, ctx context.Context) *redisStore {
	addr, _ := redis.ParseURL(redisUrl)
	rdb := redis.NewClient(addr)

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Please make sure Redis is running: %v", err)
	}

	return &redisStore{rdb}
}

func (s Shortener) store(shortUrl string, longUrl string) error {
	err := s.storage.rdb.Set(s.ctx, prefix+"::"+shortUrl, longUrl, time.Second*time.Duration(s.Timeout)).Err()
	if err != nil {
		return (err)
	}

	return nil
}

func (s Shortener) fetch(shortUrl string) (string, error) {
	longUrl, err := s.storage.rdb.Get(s.ctx, prefix+"::"+shortUrl).Result()
	if err == redis.Nil {
		err = fmt.Errorf("URL does not exist or is expired")
		return "", err
	}
	if err != nil {
		return "", err
	}

	return longUrl, nil

}

// Provide stats
func (s Shortener) Stats() *Stats {
	var stats = &Stats{}
	stats.DbSize = s.storage.rdb.DBSize(s.ctx).Val()

	return stats
}

func (s Shortener) HasCallbackUrl() bool {
	return s.CallbackUrl != ""
}
