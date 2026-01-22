package main

import (
	"crypto/rand"
	"io"
	"net/http"
	"strings"
)

var urls map[string]string

func mainPage(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, "", 400)
			return
		}

		id := rand.Text()[0:8]
		urls[id] = string(body)

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://localhost:8080/" + id))
	case http.MethodGet:
		splittedPath := strings.Split(strings.TrimPrefix(req.URL.Path, "/"), "/")
		if len(splittedPath) > 1 {
			http.Error(res, "", 400)
			return
		}

		url, ok := urls[splittedPath[0]]
		if !ok {
			http.Error(res, "", 400)
			return
		}

		res.Header().Set("Location", url)
		res.WriteHeader(307)
	}

}

func main() {
	urls = make(map[string]string)
	if err := http.ListenAndServe(":8080", http.HandlerFunc(mainPage)); err != nil {
		panic(err)
	}
}
