package main

import (
	"log"
	"net/http"

	"fga-asg-2/pkg/http/rest"
	"fga-asg-2/pkg/order"
	"fga-asg-2/pkg/storage/sqldb"
)

func main() {
	// Create storage
	// This sensitive information is written here for the convenience of this assignment
	dsn := "falfal:Pasword!2@tcp(mysql-dev-db.airy.my.id:3306)/orders_by?charset=utf8mb4&parseTime=True&loc=Local"

	storage, err := sqldb.NewStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer storage.Close()

	// Create order service
	orderService := order.NewService(storage)

	// Create HTTP server
	router := rest.NewRouter(orderService)
	http.ListenAndServe(":5000", router)

}
