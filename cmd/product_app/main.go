package main

import (
	"app/product-api/configs"
	"app/product-api/internal/rest/product"
	"app/product-api/internal/rest/product/repository"
	"app/product-api/internal/rest/user"
	db2 "app/product-api/pkg/db"
	"app/product-api/pkg/logs"
	"app/product-api/pkg/middlewares"
	"fmt"
	"log"
	"net/http"
)

func main() {
	startServer()
}

func startServer() {
	log.Println("Server is listening on port: 8080")

	logs.InitLogger()

	router := http.NewServeMux()
	conf := configs.LoadConfig()

	db, err := db2.NewDb(conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)

	// Инициализация репозитория
	productRepo := repository.NewRepository(db)
	userRepo := user.NewRepository(db)

	// Инициализация сервиса
	productService := product.NewService(productRepo)
	userService := user.NewService(userRepo)

	// Подключение хэдлеров
	product.NewHandler(router, productService)
	user.NewHandler(router, userService, conf)

	// Инициализация Middleware
	chainMdw := middlewares.CallMiddleware(
		middlewares.Cors,
		middlewares.NewAuthMiddlewares(conf.Auth.JWT).Auth,
		middlewares.Logger,
		middlewares.NewTimeoutMiddleware(600).AddTimeout,
	)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: chainMdw(router),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Server is not starting")
	}
}
