package app

import (
	"database/sql"
	"errors"
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

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS  users(\n    id SERIAL PRIMARY KEY,\n    login varchar not null unique,\n    encrypted_password varchar not null\n);\n\n\n\n\nCREATE TABLE IF NOT EXISTS bonuses   (\n    id SERIAL PRIMARY KEY,\n    user_id bigint,\n    order_id bigint,\n    change bigint,\n    type varchar(40) CHECK (type IN ('top_up', 'withdraw')),\n    status varchar(40) CHECK (status in ('NEW', 'REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED')),\n    change_date timestamp DEFAULT current_timestamp,\n    FOREIGN KEY(user_id) REFERENCES users(id)\n);\n"); err != nil {
		return nil, errors.New("error sql ")
	}

	return db, err

}
