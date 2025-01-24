package scraper

import (
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
}
