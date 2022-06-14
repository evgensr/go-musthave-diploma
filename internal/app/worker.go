package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/evgensr/go-musthave-diploma/internal/model"
	"github.com/evgensr/go-musthave-diploma/internal/store"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

type Worker struct {
	ctx    context.Context
	logger *logrus.Logger
	db     store.Store
	cfg    *Config
}

// NewWorker
func NewWorker(ctx context.Context, logger *logrus.Logger, db store.Store, cfg *Config) Worker {
	return Worker{
		ctx:    ctx,
		logger: logger,
		db:     db,
		cfg:    cfg,
	}
}

// UpdateStatus
func (w *Worker) UpdateStatus(t <-chan time.Time) {
	client := &http.Client{}
	for {
		select {
		case <-t:
			w.logger.Info("starting bonus update")
			oin := make(chan []model.Order)
			oout := make(chan model.Order)
			go w.db.User().SelectOrdersForUpdate(w.ctx, oin, oout)
			go w.getAccrual(oin, oout, client)
		case <-w.ctx.Done():
			w.logger.Info("context canceled")
		}
	}
}

func (w *Worker) getAccrual(oin chan []model.Order, oout chan model.Order, client *http.Client) {
	url := fmt.Sprintf("%s/api/orders/", w.cfg.AccrualSystemAddress)
	// w.logger.Info(url)
	orders := <-oin

	for _, order := range orders {
		var intermOrder model.AccrualOrder
		url += fmt.Sprint(order.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)

		if err != nil {
			w.logger.Fatal("request creation failed", err)
		}

		response, requestErr := w.requestWithRetry(client, request)

		if requestErr != nil {
			w.logger.Error(requestErr.Error())
		}
		defer response.Body.Close()

		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&intermOrder)
		if err != nil {
			w.logger.Debug("Error processing response" + err.Error())
		}
		log.Println("result response bonus", intermOrder)
		log.Println("bonuses: ", intermOrder.Amount)

		oout <- model.Order{ID: intermOrder.ID, Amount: intermOrder.Amount, Status: intermOrder.Status}

	}
	close(oout)
	w.logger.Info("bonus update finished")
}

// requestWithRetry
func (w *Worker) requestWithRetry(client *http.Client, request *http.Request) (*http.Response, error) {
	var response *http.Response
	var requestErr error
	for i := 0; i < 5; i++ {
		w.logger.Info(request.URL)
		response, requestErr = client.Do(request)
		if requestErr != nil {
			w.logger.Info("Retrying: " + requestErr.Error())
		} else if response.StatusCode == http.StatusTooManyRequests {
			w.logger.Info("Too many requests")
			time.Sleep(30 * time.Second)
		} else if response.StatusCode == http.StatusOK {
			return response, nil
		}
		w.logger.Info("Retrying...")
		time.Sleep(time.Duration(i*10) * time.Second)
	}

	return response, requestErr
}
