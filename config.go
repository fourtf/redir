package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
)

var (
	cfg = &config{}
)

type config struct {
	Token     string
	Addr      string
	CharCount int
}

func initConfig() {
	os.Rename("config.json", "config.json.old")

	token := randBytes()

	cfg.Addr = "http://localhost:12435"
	cfg.Token = token
	cfg.CharCount = 4

	f, err := os.Create("config.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewEncoder(f).Encode(cfg); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() {
	f, err := os.Open("config.json")
	if os.IsNotExist(err) {
		initConfig()
	} else if err != nil {
		log.Fatal("loading config.json: " + err.Error())
	}

	if err = json.NewDecoder(f).Decode(cfg); err != nil {
		log.Fatal("loading config.json: " + err.Error())
	}
}

func randBytes() string {
	b := make([]byte, 24)
	rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}
