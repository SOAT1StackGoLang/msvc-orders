package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/SOAT1StackGoLang/msvc-orders/pkg/helpers"

	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
)

// initializeApp initializes the application by loading the configuration, connecting to the datastore,
// and subscribing to the Redis channel for receiving messages.
// It returns a pointer to the RedisStore and an error if any.

var (
	binding       string
	connString    string
	paymentURI    string
	productionURI string
)

func initializeApp() (datastore.RedisStore, error) {
	flag.StringVar(&binding, "httpbind", ":8000", "address/port to bind listen socket")
	flag.Parse()
	// err := godotenv.Load()
	// if err != nil {
	// 	logger.InfoLogger.Log("load err", err.Error())
	// }
	helpers.ReadPgxConnEnvs()
	paymentURI = os.Getenv("PAYMENT_URI")
	productionURI = os.Getenv("PRODUCTION_URI")
	connString = helpers.GetConnectionParams()

	logger.InitializeLogger()

	// Load the configuration
	logger.Info("Loading configuration...")
	configs, err := LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("Connecting to datastore...")

	redisStore, err := datastore.NewRedisStore(configs.KVSURI, "", configs.KVSDB)
	if err != nil {
		// handle error
		logger.Error(err.Error())
		return nil, err
	}

	// Subscribe to the Redis channel if APP_LOG_LEVEL is set to debug
	if strings.ToLower(os.Getenv("APP_LOG_LEVEL")) == "debug" {
		err = debugChannelSubscriber(redisStore)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	return redisStore, nil
}

func debugChannelSubscriber(redisStore datastore.RedisStore) error {
	// Subscribe to the Redis channel if APP_LOG_LEVEL is set to debug
	logger.Info("DEBUG MODE ON: Subscribing to Redis channel...")
	ch, err := redisStore.SubscribeLog(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	go func() {
		for msg := range ch {
			logger.Info("channel msg: " + msg.String())
		}
	}()

	return nil
}
