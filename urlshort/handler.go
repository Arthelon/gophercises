package urlshort

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"net/http"
	"strings"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		path := strings.Trim(req.URL.EscapedPath(), " ")
		val, ok := pathsToUrls[path]
		if !ok {
			fallback.ServeHTTP(rw, req)
			return
		}
		fmt.Println(val)
		http.Redirect(rw, req, val, http.StatusSeeOther)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	records, err := parseYAML(yml)
	if err != nil {
		return fallback.ServeHTTP, err
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		path := strings.Trim(req.URL.EscapedPath(), " ")
		for _, pathRecord := range records {
			if pathRecord["path"] == path {
				http.Redirect(rw, req, pathRecord["url"], http.StatusSeeOther)
			}
		}
		fallback.ServeHTTP(rw, req)
	}, nil
}

func parseYAML(yml []byte) (dst []map[string]string, err error) {
	err = yaml.Unmarshal(yml, &dst)
	return dst, err
}
