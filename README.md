# golinky: A simple but safe URL shortener, in golang.

## Usage

### Public endpoint

`/:id:` (**get**) 

Return an HTTP redirect (307) toward the long URL.

Parameter:
- `id`: Path section of a short URL. Non existent or expired `id` will return a status code 404.

_Example:_
```bash
$ curl -i http://s.doma.in/VGtbfqWnIwVH9K0

HTTP/1.1 307 Temporary Redirect
Content-Type: text/html; charset=utf-8
Location: https://www.disney.com/abc
Date: Tue, 16 Jan 2024 16:52:06 GMT
Content-Length: 59

<a href="https://www.disney.com/">Temporary Redirect</a>.
```

### Management endpoints

>Management routes authenticate requests with the security header `x-auth-token` (or fails with status code 405).
>See `AUTH_TOKEN` above.

`/shorten` (**post**)

Return a shortened URL.

Parameters:
- `url`: long URL

_Example:_
```bash
$ curl -d '{"url": "https://www.disney.com/abc"}' \
-X POST https://s.doma.in/shorten -H "x-auth-token: mysecret"

{"shortUrl":"https://short.com/VGtbfqWnIwVH9K0"}
```

`/stats` (**get**)

Return the number of items in store.

_Example:_
```bash
$ curl https://s.doma.in:8000/stats

{"dbsize":5}
```

### Click tracking

You can optionally use a callback listener to receive click tracking data. On click events, a POST is executed toward the configured callback URL with the following data:

```json
{
    "event-type": "click",
    "msg-id": "VGtbfqWnIwVH9K0",
    "link": "https://www.disney.com/abc",
    "click-time": 1705423926,
    "user-agent": "curl/8.4.0",
    "client-ip": "127.0.0.1"
}
```

The value of `msg-id` is the path section of the shortened URL and can be used to match a sent message (email or SMS).


### Predictability of short URLs

For use cases where this is a concern, golinky mitigates the predictability of keys using a "safe" generation algorithm. You can balance security vs. usability by setting the URL length, using the `length` configuration parameter, to a value that makes sense to your use case.

## Setup 

### Configuration file

```yaml
length: 13                      # length of URL keys
timeout: 1728000                # TTL of URLs, in seconds
redis_url: ${REDIS_URL}         # URL of the Redis store
baseurl: https://s.doma.in/     # baseurl for the short URL
seed: f8J12                     # secrect seed for URL key generation
callback_url: ${CALLBACK_URL}   # callback URL for click tracking (optional)
```

### Environment variables

* `AUTH_TOKEN`: secret token to authenticate access to management endpoints (see below).


### Redis (key storage)

golinky currently requires a Redis store as a backend to keep track of shortened URLs. *Beware storage keys come with a TTL (by design)*, set using the `timeout` configuration parameter.

> Note on resilience: since Redis isn't natively resilient to crashes and restart, make sure you implement some form of data persistence mechanism if you need any level of resilience (backups of else).
