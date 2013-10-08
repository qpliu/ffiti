package app

import (
	"ffiti"
	"net/http"
)

func init() {
	if err := ffiti.Init(http.DefaultServeMux, ffiti.Config{
		DataStore: DataStore(true),
		Documents: "documents/",
	}); err != nil {
		panic(err.Error())
	}
}
