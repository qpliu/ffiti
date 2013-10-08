package standalone

// For simple testing only.  Assumes single-threading, so
// concurrent modifications will result in corrupt data.

import (
	"ffiti"
	"ffiti/data"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	MAX_POSTS_PER_KEY = 100
)

type DataStore string

func NewDataStore(root string) (DataStore, error) {
	if err := os.MkdirAll(root, os.ModePerm); err != nil {
		return DataStore(root), err
	}
	return DataStore(root), nil
}

func (ds DataStore) Storage(r *http.Request) ffiti.Storage {
	return ds
}

func (ds DataStore) Post(key string, post data.Post) error {
	file, err := os.OpenFile(filepath.Join(string(ds), key), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	posts, err := data.ReadPosts(file)
	if err != nil && err != io.EOF {
		return err
	}
	if _, err = file.Seek(0, os.SEEK_SET); err != nil {
		return err
	}
	if err = file.Truncate(0); err != nil {
		return err
	}
	return data.AddPost(post, posts, file, MAX_POSTS_PER_KEY)
}

func (ds DataStore) Get(key string) ([]data.Post, error) {
	file, err := os.Open(filepath.Join(string(ds), key))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()
	return data.ReadPosts(file)
}
