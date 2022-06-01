package app

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evgensr/go-musthave-diploma/internal/helper"
	"github.com/evgensr/go-musthave-diploma/internal/store/pg"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

// APIserver ...
type APIserver struct {
	config       *Config
	logger       *logrus.Logger
	router       *mux.Router
	store        Storage
	sessionStore sessions.Store
}

// New ...
func New(config *Config, sessionStore sessions.Store) *APIserver {
	var store Storage
	// sessionStore = sessions.NewCookieStore([]byte(config.SessionKey))
	config.BaseURL = helper.AddSlash(config.BaseURL)
	// param := config.FileStoragePath
	store = pg.New(config.DatabaseURI)

	return &APIserver{
		config:       config,
		logger:       logrus.New(),
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}
}

// Start ...
func (s *APIserver) Start() error {

	if err := s.configureLogger(); err != nil {
		return err
	}

	if len(s.config.DatabaseURI) > 1 {
		if err := s.CreateTable(); err != nil {
			log.Fatal("create table ", err)
		}
	}

	s.configureRouter()

	s.logger.Info("SERVER_ADDRESS ", s.config.RunAddress)
	s.logger.Info("BASE_URL ", s.config.BaseURL)

	s.logger.Info("Staring api server")
	return http.ListenAndServe(s.config.RunAddress, s.router)
}

func (s *APIserver) configureLogger() error {

	level, err := logrus.ParseLevel(s.config.LogLevel)

	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil

}

func (s *APIserver) configureRouter() {

	s.router.Use(s.authenticateUser)
	s.router.HandleFunc("/ping", s.HandlerPing())

	// регистрация пользователя;
	s.router.HandleFunc("/api/user/register", s.HandlerUserRegister()).Methods(http.MethodPost)

	// аутентификация пользователя;
	s.router.HandleFunc("/api/user/login", s.HandlerUserLogin()).Methods(http.MethodPost)

	// загрузка пользователем номера заказа для расчёта;
	s.router.HandleFunc("/api/user/orders", s.HandlerUserOrders()).Methods(http.MethodPost)

	// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
	s.router.HandleFunc("/api/user/orders", s.HandlerUserGetOrders()).Methods(http.MethodGet)

	// получение текущего баланса счёта баллов лояльности пользователя;
	s.router.HandleFunc("/api/user/balance", s.HandlerUserGetBalance()).Methods(http.MethodGet)

	// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	s.router.HandleFunc("/api/user/balance/withdraw", s.HandlerUserBalanceWithdraw()).Methods(http.MethodPost)

	// получение информации о выводе средств с накопительного счёта пользователем.
	s.router.HandleFunc("/api/user/balance/withdrawals", s.HandlerUserGetBalanceWithdrawals()).Methods(http.MethodGet)

	// взаимодействие с системой расчёта начислений баллов лояльности
	s.router.HandleFunc("/api/orders/{number}", s.HandlerUserGeOrders()).Methods(http.MethodGet)

	s.router.Use(s.GzipHandleEncode)
	s.router.Use(s.GzipHandleDecode)

}

func (s *APIserver) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var id string

		log.Println("SessionKey:  ", s.config.SessionKey)

		c, err := r.Cookie(sessionName)
		if err != nil {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			id = helper.GeneratorUUID()
			encryptedCookie, err := helper.Encrypted([]byte(id), s.config.SessionKey)
			if err != nil {
				s.logger.Warning("error encrypted cookie ", err)
				return
			}
			cookie := http.Cookie{Name: sessionName, Value: hex.EncodeToString(encryptedCookie), Expires: expiration}
			http.SetCookie(w, &cookie)
		} else {
			fmt.Println("Cookie ", c.Value)
			decoded, err := hex.DecodeString(c.Value)
			if err != nil {
				s.logger.Warning("error decode string Cookie ", c.Value)
				return
			}
			decryptedCookie, err := helper.Decrypted(decoded, s.config.SessionKey)
			if err != nil {
				// не смогли декодировать, устанавливаем новую куку и юзера
				expiration := time.Now().Add(365 * 24 * time.Hour)
				id = helper.GeneratorUUID()
				encryptedCookie, err := helper.Encrypted([]byte(id), s.config.SessionKey)
				if err != nil {
					s.logger.Warning("error encrypted cookie ", err)
					return
				}
				cookie := http.Cookie{Name: sessionName, Value: hex.EncodeToString(encryptedCookie), Expires: expiration}
				http.SetCookie(w, &cookie)

			}

			id = string(decryptedCookie)
		}

		log.Println("user id: ", id)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, id)))

	})

}

func (s *APIserver) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

func (s *APIserver) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *APIserver) CreateTable() error {
	log.Println("config Database: ", s.config.DatabaseURI)
	db, err := sql.Open("postgres", s.config.DatabaseURI)
	if err != nil {
		log.Println("create table func ", err)
		return err
	}
	if err := db.Ping(); err != nil {
		log.Println("ping err ", err)
		return err
	}

	if _, err := db.Exec("CREATE TABLE  IF NOT EXISTS short" +
		"(id serial primary key," +
		"original_url varchar(4096) not null," +
		"short_url varchar(32) UNIQUE not null," +
		"user_id varchar(36) not null," +
		"correlation_id varchar(36) null," +
		"status smallint not null DEFAULT 0);"); err != nil {
		return errors.New("error sql ")
	}

	return nil

}
