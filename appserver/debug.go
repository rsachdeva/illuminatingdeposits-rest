package appserver

import (
	"log"
	"net/http"
)

func Debug(log *log.Logger, cfg AppConfig) {
	log.Println("debug responder listening on", cfg.Web.Debug)
	err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
	log.Println("debug responder closed", err)
}
