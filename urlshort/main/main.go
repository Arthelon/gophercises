package main

import (
	"flag"
	"fmt"
	"github.com/arthelon/gophercises/urlshort"
	"io/ioutil"
	"net/http"
)

var (
	yamlPath string
)

func init() {
	flag.StringVar(&yamlPath, "yaml", "", "File to parse for url shortener")
	flag.Parse()
}


func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var yaml string
	if len(yamlPath) > 0 {
		b, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			panic(err)
		}
		yaml = string(b)
	} else {
		yaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/final
`
	}
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
