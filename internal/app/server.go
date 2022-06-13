package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/evgensr/go-musthave-diploma/internal/model"
	"github.com/evgensr/go-musthave-diploma/internal/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/theplant/luhn"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	sessionName        = "education"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect login or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
	ctx          context.Context
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}

}

func (s *server) handlerPostOrders() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// создаем переменную для передачи в функции
		order := model.Order{}

		orderID, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, struct {
				Status string `json:"status"`
			}{Status: "invalid request format"})
			return
		}

		// log.Println(string(orderID))
		// log.Println(r.Context().Value(ctxKeyUser))

		// id, _ := strconv.Atoi(string(orderID))
		id, _ := strconv.Atoi(string(orderID))

		// проверка луна
		ok := luhn.Valid(id)
		if !ok {
			s.respond(w, r, 422, struct {
				Status string `json:"status"`
			}{Status: "invalid order number format"})
			return
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)

		order.UserID = user.ID
		order.Status = "NEW"
		order.Type = "top_up"
		order.ID = string(orderID)

		expectedUser, err := s.store.User().SelectUserForOrder(s.ctx, order)
		if expectedUser == user.ID {
			s.respond(w, r, 200, struct {
				Status string `json:"status"`
			}{Status: "the order number has already been uploaded by this user"})
			return
		}
		if expectedUser != 0 {
			s.respond(w, r, 409, struct {
				Status string `json:"status"`
			}{Status: "the order number has already been uploaded by other user"})
			return
		}

		err = s.store.User().InsertOrder(s.ctx, order)

		s.respond(w, r, http.StatusAccepted, struct {
			Status string `json:"status"`
		}{Status: "accepted"})

	}

}

func (s *server) handlerGetOrders() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(ctxKeyUser).(*model.User)

		orders, err := s.store.User().SelectAllOrders(s.ctx, user.ID)

		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, struct {
				Status string `json:"status"`
			}{Status: err.Error()})
			return
		}

		if len(orders) == 0 {
			s.respond(w, r, 204, struct {
				Status string `json:"status"`
			}{Status: "no data for the user"})
			return
		}

		s.respond(w, r, http.StatusOK, orders)

	}

}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Println(r.Body)
			s.error(w, r, http.StatusBadRequest, err)
			return

		}

		u := &model.User{
			Login:    req.Login,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, 409, err)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		s.sessionStore.Save(r, w, session)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusOK, u)

	}

}

func (s *server) handlerGetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(ctxKeyUser).(*model.User)

		orders, err := s.store.User().SelectBalance(s.ctx, user.ID)

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, orders)

	}

}

// handlerPostWithdraw  списания бонусов /api/user/balance/withdraw
func (s *server) handlerPostWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var o model.Withdrawal
		err := decoder.Decode(&o)

		// получаем текущего пользователя из контекста
		user := r.Context().Value(ctxKeyUser).(*model.User)

		// преобразуем строку ордера в число int для проверки луна
		i, err := strconv.Atoi(o.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// проверка луна
		ok := luhn.Valid(i)
		if !ok {
			s.error(w, r, 422, errors.New("invalid order"))
			return
		}

		order := model.Order{ID: o.ID, Amount: -o.Amount, UserID: user.ID, Status: "PROCESSED", Type: "withdraw"}
		err = s.store.User().InsertOrder(s.ctx, order)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, struct {
			Status string `json:"status"`
		}{Status: "ok"})

	}

}

// handlerGetWithdraw Получение информации о выводе средств /api/user/balance/withdrawals
func (s *server) handlerGetWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// получаем текущего пользователя из контекста
		user := r.Context().Value(ctxKeyUser).(*model.User)
		// получаем количество бонусов из бд
		result, err := s.store.User().SelectAllWithdrawals(s.ctx, user.ID)
		// если ошибка, то выводим 500
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		// проверяем существование записей о списании
		//if len(result) == 0 {
		//	s.error(w, r, 204, errors.New("result empty"))
		//	return
		//}

		s.respond(w, r, http.StatusAccepted, result)

	}

}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
