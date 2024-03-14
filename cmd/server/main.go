// just a simple example buildable
package main

import (
	"log"
	"net/http"

	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/transport"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/transport/routes"
	paymentsapi "github.com/SOAT1StackGoLang/msvc-payments/pkg/api"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	productionapi "github.com/SOAT1StackGoLang/msvc-production/pkg/api"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cache, err := initializeApp()
	if err != nil {
		panic("unable to connect")
	}

	gormDB, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Panicf("failed initializing db: %s\n", err)
	}

	r := mux.NewRouter()

	catRepo := persistence.NewCategoriesPersistence(gormDB, logger.InfoLogger)
	categoriesSvc := service.NewCategoriesService(catRepo, logger.InfoLogger)
	r = routes.NewCategoriesRouter(categoriesSvc, r, logger.InfoLogger)

	productsRepo := persistence.NewProductsPersistence(gormDB, logger.InfoLogger)
	productsSvc := service.NewProductsService(productsRepo, logger.InfoLogger)
	r = routes.NewProductsRouter(productsSvc, r, logger.InfoLogger)

	productionClient := productionapi.NewClient(productionURI, logger.InfoLogger)
	paymentsClient := paymentsapi.NewClient(paymentURI, logger.InfoLogger)

	paymentsRepo := persistence.NewPaymentsPersistence(gormDB, logger.InfoLogger)
	paymentsSvc := service.NewPaymentsService(paymentsRepo, paymentsClient, logger.InfoLogger, cache)
	r = routes.NewPaymentsRouter(paymentsSvc, r, logger.InfoLogger)

	ordersRepo := persistence.NewOrdersPersistence(gormDB, logger.InfoLogger)
	ordersSvc := service.NewOrdersService(ordersRepo, productsSvc, paymentsSvc, logger.InfoLogger, cache, paymentsClient, productionClient)
	r = routes.NewOrdersRouter(ordersSvc, r, logger.InfoLogger)

	transport.NewHTTPServer(":8080", muxToHttp(r))
}

func muxToHttp(r *mux.Router) http.Handler {
	return r
}
