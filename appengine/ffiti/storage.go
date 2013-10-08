package ffiti

import (
	"ffiti/data"
	"net/http"
)

type DataStore interface {
	Storage(r *http.Request) Storage
}

type Storage interface {
	Post(key string, post data.Post) error
	Get(key string) ([]data.Post, error)
}
