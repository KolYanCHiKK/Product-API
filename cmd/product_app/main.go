package main

import (
	"app/product-api/configs"
	"app/product-api/internal/rest/product"
	db2 "app/product-api/pkg/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	startServer()
}

func startServer() {
	log.Println("Server is listening on port: 8080")

	router := http.NewServeMux()
	conf := configs.LoadConfig()

	db, err := db2.NewDb(conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)

	product.NewHandler(router, conf)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Server is not starting")
	}
}
