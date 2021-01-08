package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

type AppConfig struct {
	Web struct {
		Address         string        `conf:"default:0.0.0.0:3000"`
		Debug           string        `conf:"default:0.0.0.0:4000"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	DB struct {
		User       string `conf:"default:postgres"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:db"`
		Name       string `conf:"default:postgres"`
		DisableTLS bool   `conf:"default:true"`
	}
	Trace struct {
		URL         string  `conf:"default:http://zipkin:9411/api/v2/spans"`
		Service     string  `conf:"default:illuminatingdeposits-rest"`
		Probability float64 `conf:"default:1"`
	}
}

func ParsedConfig(cfg AppConfig) (AppConfig, error) {
	if err := conf.Parse(os.Args[1:], "DEPOSITS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("DEPOSITS", &cfg)
			if err != nil {
				return AppConfig{}, errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return AppConfig{}, nil
		}
		return AppConfig{}, errors.Wrap(err, "parsing config")
	}
	log.Printf("AppConfig is %v", cfg)
	return cfg, nil
}

func tlsConfig() (*tls.Config, error) {
	certFile := "conf/tls/servercrtto.pem"
	keyFile := "conf/tls/serverkeyto.pem"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "LoadX509KeyPair error")
	}
	fmt.Println("No errors with LoadX509KeyPair")
	tl := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return &tl, nil
}
