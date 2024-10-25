module github.com/dbyio/golinky

go 1.23

require (
	github.com/gorilla/mux v1.8.1
	internal/callback v0.0.0-00010101000000-000000000000
	internal/urlshortener v1.0.0
)

replace internal/urlshortener => ./internal/urlshortener

replace internal/callback => ./internal/callback

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	golang.org/x/sys v0.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
