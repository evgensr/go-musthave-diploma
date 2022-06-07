package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/evgensr/go-musthave-diploma/internal/app"
	"log"
	"os"
)

var version = "0.0.1"

func main() {

	conf := app.NewConfig()
	err := env.Parse(&conf)

	if err != nil {
		log.Printf("[ERROR] failed to parse flags: %v", err)
		os.Exit(1)
	}

	flag.StringVar(&conf.RunAddress, "a", conf.RunAddress, "RUN_ADDRESS")
	flag.StringVar(&conf.BaseURL, "b", conf.BaseURL, "BASE_URL")
	flag.StringVar(&conf.DatabaseURI, "d", conf.DatabaseURI, "DATABASE_URI")
	flag.StringVar(&conf.AccrualSystemAddress, "r", conf.AccrualSystemAddress, "ACCRUAL_SYSTEM_ADDRESS")

	flag.Parse()

	// sessionStore := sessions.NewCookieStore([]byte(conf.SessionKey))

	if err := app.Start(&conf); err != nil {
		log.Fatal("fatal ", err)
	}

	//
	//server := app.New(&conf, sessionStore)
	//
	//if err := server.Start(); err != nil {
	//	log.Fatal(err)
	//}

}
