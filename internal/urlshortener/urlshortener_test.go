package urlshortener

import (
	"context"
	"strings"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

var testShortener *Shortener

const longUrl_1 = "https://this.is.a/veryveryverylongurl"

func init() {
	length := uint8(15)
	timeout := 3
	redissrv := "127.0.0.1:6379"
	redispassword := ""
	baseurl := "https://short.com/"
	seed := "test"
	callbackurl := ""

	testShortener = &Shortener{
		length,
		timeout,
		redissrv,
		redispassword,
		baseurl,
		seed,
		callbackurl,
		nil,
		nil,
	}
	(*testShortener).ctx = context.Background()
	(*testShortener).storage = initStore(testShortener.redisSrv,
		testShortener.redisPassword, testShortener.ctx)
}

func TestCreateShortUrl(t *testing.T) {
	url_1, err := testShortener.CreateShortUrl(longUrl_1)
	assert.NilError(t, err)
	time.Sleep(10 * time.Millisecond)
	url_2, err := testShortener.CreateShortUrl(longUrl_1)
	assert.NilError(t, err)

	assert.Assert(t, url_1 != url_2)
}

func TestUseShortUrl(t *testing.T) {
	url_1, err := testShortener.CreateShortUrl(longUrl_1)
	assert.NilError(t, err)

	id := strings.ReplaceAll(url_1, testShortener.baseUrl, "")

	url_2, err := testShortener.UseShortUrl(id)
	assert.NilError(t, err)
	assert.Equal(t, url_2, longUrl_1)
}

func TestItemExpires(t *testing.T) {
	url_1, err := testShortener.CreateShortUrl(longUrl_1)
	assert.NilError(t, err)
	id_1 := strings.ReplaceAll(url_1, testShortener.baseUrl, "")

	url, _ := testShortener.UseShortUrl(id_1)
	assert.Equal(t, url, longUrl_1)

	time.Sleep(time.Duration(testShortener.timeout+1) * time.Second)
	url, _ = testShortener.UseShortUrl(id_1)
	assert.Assert(t, url == "")
}
