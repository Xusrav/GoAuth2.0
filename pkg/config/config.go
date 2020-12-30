package config

import (
	"log"
	"os"

	"github.com/creamdog/gonfig"
)

var (
	ClientID, ClientSecret, Host, Port, URLomdbApi, ApiKey string
)

func init() {

	f, err := os.Open("../configs/config.json")
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	config, err := gonfig.FromJson(f)
	if err != nil {
		log.Print(err)
	}

	ClientID, err = config.GetString("client_id", "")
	if err != nil {
		log.Print(err)
	}
	ClientSecret, err = config.GetString("client_secret", "")
	if err != nil {
		log.Print(err)
	}
	Host, err = config.GetString("host", "localhost")
	if err != nil {
		log.Print(err)
	}
	Port, err = config.GetString("port", "9999")
	if err != nil {
		log.Print(err)
	}

	URLomdbApi, err = config.GetString("url_omdb", "")
	if err != nil {
		log.Print(err)
	}

	ApiKey, err = config.GetString("api_key", "")
	if err != nil {
		log.Print(err)
	}
}
