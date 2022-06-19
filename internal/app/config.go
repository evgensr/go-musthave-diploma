package app

// Config ...
type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"0.0.0.0:8080"`
	BaseURL              string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"debug"`
	SessionKey           string `env:"SESSION_KEY" envDefault:"SESSION_KEY"`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:"host=localhost user=postgres password=postgres dbname=restapi sslmode=disable"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// NewConfig ...
func NewConfig() Config {
	return Config{}
}
