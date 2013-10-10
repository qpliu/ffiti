package app

import (
	"ffiti"
	"io/ioutil"
	"net/http"
	"regexp"
)

func init() {
	if err := ffiti.Init(http.DefaultServeMux, ffiti.Config{
		DataStore: DataStore(true),
		Documents: "documents/",
		Version:   version(),
	}); err != nil {
		panic(err.Error())
	}
}

func version() string {
	file, err := ioutil.ReadFile("app.yaml")
	if err != nil {
		panic(err.Error())
	}
	return regexp.MustCompile(".*\nversion: *(.*)\n").FindStringSubmatch(string(file))[1]
}
