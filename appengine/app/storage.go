package app

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"ffiti"
	"ffiti/data"
	"io"
	"net/http"
)

const (
	MAX_POSTS_PER_KEY = 100
)

type DataStore bool

func (ds DataStore) Storage(r *http.Request) ffiti.Storage {
	return Storage{context: appengine.NewContext(r)}
}

type Storage struct {
	context appengine.Context
}

func makeKey(context appengine.Context, key string) *datastore.Key {
	return datastore.NewKey(context, "ffiti", key, 0, nil)
}

type record struct {
	Bytes []byte
}

func (s Storage) Post(key string, post data.Post) error {
	var status error
	if err := datastore.RunInTransaction(s.context, func(context appengine.Context) error {
		k := makeKey(context, key)
		value, err := get(context, k)
		if err != nil {
			return err
		}
		posts, status := data.ReadPosts(bytes.NewBuffer(value.Bytes))
		if status != nil && status != io.EOF {
			return nil
		}
		buffer := bytes.NewBuffer(nil)
		status = data.AddPost(post, posts, buffer, MAX_POSTS_PER_KEY)
		if status != nil {
			return nil
		}
		_, err = datastore.Put(context, k, &record{Bytes: buffer.Bytes()})
		return err
	}, nil); err != nil {
		return err
	}
	return status
}

func (s Storage) Get(key string) ([]data.Post, error) {
	k := makeKey(s.context, key)
	value, err := get(s.context, k)
	if err != nil {
		return nil, err
	}
	return data.ReadPosts(bytes.NewBuffer(value.Bytes))
}

func get(context appengine.Context, k *datastore.Key) (record, error) {
	value := record{}
	if err := datastore.Get(context, k, &value); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return value, nil
		}
		return value, err
	}
	return value, nil
}
