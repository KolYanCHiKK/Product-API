package main

import (
	"app/product-api/configs"
	"app/product-api/internal/rest/order"
	"app/product-api/internal/rest/product/repository"
	"app/product-api/internal/rest/user"
	db2 "app/product-api/pkg/db"
)

func main() {
	conf := configs.LoadConfig()
	db, err := db2.NewDb(conf)
	if err != nil {
		panic("ERROR DB CONNECT")
	}

	err = db.AutoMigrate(
		repository.Product{}, user.User{}, user.Session{},
		order.Order{}, order.OrderProduct{},
	)
	if err != nil {
		panic("MIGRATION FAILED")
	}
}
