package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// HandlerUserRegister - Регистрация пользователя
func (s *APIserver) HandlerUserRegister() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var request request
		var err error

		// создаем перменную
		decoder := json.NewDecoder(r.Body)

		// декодируем в структуру request
		err = decoder.Decode(&request)
		if err != nil {
			// логируем ошибку
			s.logger.Warning("HandlerSetURL: request not json ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(request)
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerUserLogin - Аутентификация пользователя
func (s *APIserver) HandlerUserLogin() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var request request
		var err error

		// создаем перменную
		decoder := json.NewDecoder(r.Body)

		// декодируем в структуру request
		err = decoder.Decode(&request)
		if err != nil {
			// логируем ошибку
			s.logger.Warning("HandlerSetURL: request not json ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(request)
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerUserOrders - Загрузка номера заказа
func (s *APIserver) HandlerUserOrders() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var request request
		var err error

		// создаем перменную
		decoder := json.NewDecoder(r.Body)

		// декодируем в структуру request
		err = decoder.Decode(&request)
		if err != nil {
			// логируем ошибку
			s.logger.Warning("HandlerSetURL: request not json ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(request)
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerUserGetOrders - Получение списка загруженных номеров заказов
func (s *APIserver) HandlerUserGetOrders() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var request request
		var err error

		// создаем перменную
		decoder := json.NewDecoder(r.Body)

		// декодируем в структуру request
		err = decoder.Decode(&request)
		if err != nil {
			// логируем ошибку
			s.logger.Warning("HandlerSetURL: request not json ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(request)
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerUserGetBalance - Получение текущего баланса пользователя
func (s *APIserver) HandlerUserGetBalance() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var request request
		var err error

		// создаем перменную
		decoder := json.NewDecoder(r.Body)

		// декодируем в структуру request
		err = decoder.Decode(&request)
		if err != nil {
			// логируем ошибку
			s.logger.Warning("HandlerSetURL: request not json ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(request)
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerUserBalanceWithdraw - Запрос на списание средств
func (s *APIserver) HandlerUserBalanceWithdraw() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var request request
		var err error

		// создаем перменную
		decoder := json.NewDecoder(r.Body)

		// декодируем в структуру request
		err = decoder.Decode(&request)
		if err != nil {
			// логируем ошибку
			s.logger.Warning("HandlerSetURL: request not json ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(request)
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerUserGetBalanceWithdrawals - Получение информации о выводе средств
func (s *APIserver) HandlerUserGetBalanceWithdrawals() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var request request
		var err error

		// создаем перменную
		decoder := json.NewDecoder(r.Body)

		// декодируем в структуру request
		err = decoder.Decode(&request)
		if err != nil {
			// логируем ошибку
			s.logger.Warning("HandlerSetURL: request not json ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(request)
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerUserGeOrders - Взаимодействие с системой расчёта начислений баллов лояльности
func (s *APIserver) HandlerUserGeOrders() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var request request
		var err error

		// создаем перменную
		decoder := json.NewDecoder(r.Body)

		// декодируем в структуру request
		err = decoder.Decode(&request)
		if err != nil {
			// логируем ошибку
			s.logger.Warning("HandlerSetURL: request not json ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(request)
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerPing проверка подключения к базе данных
func (s *APIserver) HandlerPing() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// s.config.DatabaseDSN = "host=localhost user=postgres password=postgres dbname=restapi sslmode=disable"
		log.Println("DatabaseDSN: ", s.config.DatabaseURI)
		db, err := sql.Open("postgres", s.config.DatabaseURI)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := db.Ping(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	}

}
