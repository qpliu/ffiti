package ffiti

import (
	"encoding/json"
	"ffiti/data"
	"log"
	"net/http"
)

const (
	VERSION = "0.0"
)

type Config struct {
	DataStore DataStore
	Documents string
}

type Posts struct {
	Bounds [4]float64  `json:"bounds"`
	Posts  []data.Post `json:"posts"`
}

func Init(serveMux *http.ServeMux, config Config) error {
	serveMux.Handle("/", http.FileServer(http.Dir(config.Documents)))

	serveMux.HandleFunc("/v1/post", func(w http.ResponseWriter, r *http.Request) {
		post, err := data.GetPost(r)
		if err != nil {
			log.Printf("data.GetPost,err=%s", err.Error())
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		err = config.DataStore.Storage(r).Post(post.Location.Key(), *post)
		if err != nil {
			log.Printf("config.Storage.GetPost,err=%s", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Printf("/v1/post:loc=%v",post.Location)
		w.WriteHeader(http.StatusCreated)
	})

	serveMux.HandleFunc("/v1/get", func(w http.ResponseWriter, r *http.Request) {
		loc, err := data.GetLocation(r)
		if err != nil {
			log.Printf("data.GetLocation,err=%s", err.Error())
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		log.Printf("/v1/get:loc=%v",loc)
		result := Posts{Bounds: loc.Bounds()}
		storage := config.DataStore.Storage(r)
		for _, key := range loc.Keys() {
			if msgs, err := storage.Get(key); err == nil {
				for _, msg := range msgs {
					result.Posts = append(result.Posts, msg)
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	serveMux.HandleFunc("/v1/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(VERSION))
	})

	return nil
}
