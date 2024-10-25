package urlshortener

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

func hashMe(data string, seed string) string {
	ts := fmt.Sprintf("%d", time.Now().UnixMilli())
	blob := []byte(data + seed + ts)

	h := sha256.New()
	h.Write(blob)
	ret := encodeBase62(h.Sum(nil))
	return ret
}

func encodeBase62(data []byte) string {
	var i big.Int
	i.SetBytes(data[:])
	return i.Text(62)
}

func loadConfig(file string) *Shortener {
	type config struct {
		Length      int    `yaml:"length"`
		Timeout     int    `yaml:"timeout"`
		RedisUrl    string `yaml:"redis_url"`
		BaseUrl     string `yaml:"baseurl"`
		Seed        string `yaml:"seed"`
		CallbackUrl string `yaml:"callback_url"`
	}

	f, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var c config

	if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(f))), &c); err != nil {
		log.Fatal(err)
	}

	s := &Shortener{
		uint8(c.Length), c.Timeout, c.RedisUrl,
		c.BaseUrl, c.Seed, c.CallbackUrl, nil, nil,
	}
	return s
}
