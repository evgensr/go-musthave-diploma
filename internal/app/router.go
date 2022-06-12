package app

import (
	"github.com/gorilla/handlers"
	"net/http"
)

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/api/user/register", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/api/user/login", s.handleSessionsCreate()).Methods("POST")

	private := s.router.PathPrefix("").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/api/user/orders", s.handlerPostOrders()).Methods(http.MethodPost)
	private.HandleFunc("/api/user/orders", s.handlerGetOrders()).Methods(http.MethodGet)
	private.HandleFunc("/api/user/balance", s.handlerGetBalance()).Methods(http.MethodGet)
	private.HandleFunc("/api/user/balance/withdraw", s.handlerPostWithdraw()).Methods(http.MethodPost)
	private.HandleFunc("/api/user/balance/withdrawals", s.handlerGetWithdraw()).Methods(http.MethodGet)

	// s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	// регистрация пользователя;
	// s.router.HandleFunc("/api/user/register", s.HandlerUserRegister()).Methods(http.MethodPost)

	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")

}
