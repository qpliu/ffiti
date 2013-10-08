package main

import (
	"./standalone"
	"ffiti"
	"net/http"
	"os"
)

func main() {
	documents := "./appengine/documents"
	data := "/tmp/ffiti"
	if len(os.Args) > 1 {
		documents = os.Args[1]
	}
	serveMux := http.NewServeMux()
	dataStore, err := standalone.NewDataStore(data)
	if err != nil {
		panic(err.Error())
	}
	if err := ffiti.Init(serveMux, ffiti.Config{
		DataStore: dataStore,
		Documents: documents,
	}); err != nil {
		panic(err.Error())
	}
	if err := http.ListenAndServe(":8080", serveMux); err != nil {
		panic(err.Error())
	}
}
