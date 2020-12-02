package rest

import (
	"log"
	"net/http"
)

func Debug(log *log.Logger, cfg AppConfig) {
	log.Println("debug router listening on", cfg.Web.Debug)
	err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
	log.Println("debug router closed", err)
}
