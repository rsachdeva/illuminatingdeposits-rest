package rest

import (
	"log"
	"net/http"
)

func Debug(log *log.Logger, cfg AppConfig) {
	log.Println("debug service listening on", cfg.Web.Debug)
	err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
	log.Println("debug service closed", err)
}
