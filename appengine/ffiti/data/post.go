package data

import (
	"encoding/gob"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	MAX_MESSAGE_LENGTH = 140
)

type Post struct {
	Location  Location  `json:"loc"`
	Message   string    `json:"msg"`
	Timestamp time.Time `json:"t"`
}

func GetPost(r *http.Request) (*Post, error) {
	loc, err := GetLocation(r)
	if err != nil {
		return nil, err
	}
	msg := r.FormValue("msg")
	if msg == "" {
		return nil, errors.New("Empty msg")
	}
	if len(msg) > MAX_MESSAGE_LENGTH {
		msg = msg[:MAX_MESSAGE_LENGTH]
	}
	return &Post{
		Location:  *loc,
		Message:   msg,
		Timestamp: time.Now(),
	}, nil
}

func init() {
	gob.Register(Post{Location: Location{}})
}

func ReadPosts(r io.Reader) ([]Post, error) {
	var count int
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&count); err != nil {
		return nil, err
	}
	posts := make([]Post, count)
	for i := 0; i < count; i++ {
		if err := decoder.Decode(&posts[i]); err != nil {
			return posts, err
		}
	}
	return posts, nil
}

func AddPost(post Post, posts []Post, w io.Writer, maxPosts int) error {
	replace := -1
	encoder := gob.NewEncoder(w)
	count := len(posts)
	if count >= maxPosts {
		if err := encoder.Encode(count); err != nil {
			return err
		}
		replace = 0
		for i := 1; i < count; i++ {
			if posts[i].Timestamp.Before(posts[replace].Timestamp) {
				replace = i
			}
		}
	} else {
		if err := encoder.Encode(count + 1); err != nil {
			return err
		}
		if err := encoder.Encode(post); err != nil {
			return err
		}
	}
	for i := 0; i < count; i++ {
		if i == replace {
			if err := encoder.Encode(post); err != nil {
				return err
			}
		} else {
			if err := encoder.Encode(posts[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
