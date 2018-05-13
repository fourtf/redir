package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

var (
	letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func main() {
	os.Mkdir("/var/redir", os.ModePerm)

	// load config
	loadConfig()

	// create router
	r := chi.NewRouter()

	r.Post("/r/add", handleAdd)
	r.Get("/r/{id:[0-9a-zA-Z]+}", handleGet)

	// create server
	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	srv.ListenAndServe()
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	id := randSeq(cfg.CharCount)

	if r.Header.Get("Token") != cfg.Token {
		http.Error(w, "401 Unauthorized.", http.StatusUnauthorized)
		return
	}

	if r.Header.Get("Url") == "" {
		http.Error(w, "400 Bad request.", http.StatusBadRequest)
		return
	}

	f, err := os.Create("/var/redir/" + id)

	if err != nil {
		http.Error(w, "502 Internal server error.", 502)
		return
	}

	f.WriteString(r.Header.Get("Url"))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	f, err := os.Open("/var/redir/" + id)
	if err != nil {
		http.Error(w, "404 Not found.", 404)
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, "502 Internal server error.", 502)
		return
	}

	http.Redirect(w, r, string(b), http.StatusTemporaryRedirect)
}

// source: https://stackoverflow.com/a/22892986
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
