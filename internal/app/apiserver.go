package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/evgensr/go-musthave-diploma/internal/store/sqlstore"
	"github.com/gorilla/sessions"
)

func Start(config *Config) error {

	log.Println(config)
	db, err := newDB(config.DatabaseURI)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)
	return http.ListenAndServe(config.RunAddress, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, err

}
