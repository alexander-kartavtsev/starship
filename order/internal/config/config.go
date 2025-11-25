package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/alexander-kartavtsev/starship/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
	Postgres      PostgresConfig
	Server        ServerHTTPConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryGRPCConfig, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCConfig, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	serverHttpConfig, err := env.NewServerHttpConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerConfig,
		InventoryGRPC: inventoryGRPCConfig,
		PaymentGRPC:   paymentGRPCConfig,
		Postgres:      postgresConfig,
		Server:        serverHttpConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
