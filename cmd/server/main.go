// just a simple example buildable
package main

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/persistence"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/transport"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	_, err := initializeApp()
	log.Println("Bootstrapping msvc-orders...")

	ctx := context.Background()

	gormDB, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Panicf("failed initializing db: %s\n", err)
	}

	catRepo := persistence.NewCategoriesPersistence(ctx, gormDB, logger.InfoLogger)
	categoriesSvc := service.NewCategoriesService(catRepo, logger.InfoLogger)

	handler := transport.NewHTTPHandler(categoriesSvc, logger.InfoLogger)

	logger.Info("Starting http server...")
	transport.NewHTTPServer(":8080", handler)
}
