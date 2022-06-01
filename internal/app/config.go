package app

import (
	"github.com/evgensr/go-musthave-diploma/internal/store"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS,required" envDefault:"0.0.0.0:8080"`
	BaseURL              string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"debug"`
	SessionKey           string `env:"SESSION_KEY" envDefault:"SESSION_KEY"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() Config {
	return Config{}
}

type key int

const (
	sessionName     = "practicum"
	ctxKeyUser  key = iota
)

type request struct {
	URL string `json:"url" valid:"url"`
}

type response struct {
	URL string `json:"result" valid:"url"`
}

type Line = store.Line

type Storage interface {
	Get(key string) (Line, error)
	Set(Line) error
	Delete([]Line) error
	GetByUser(key string) []Line
}
