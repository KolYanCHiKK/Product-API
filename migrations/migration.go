package main

import (
	"app/product-api/configs"
	"app/product-api/internal/rest/product"
	db2 "app/product-api/pkg/db"
)

func main() {
	conf := configs.LoadConfig()
	db, err := db2.NewDb(conf)
	if err != nil {
		panic("ERROR DB CONNECT")
	}

	err = db.AutoMigrate(product.Product{})
	if err != nil {
		panic("MIGRATION FAILED")
	}
}
