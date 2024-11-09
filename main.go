package main

import (
	"github.com/ricky2122/go-echo-example/infrastructure"
	"github.com/ricky2122/go-echo-example/infrastructure/api"
)

func main() {
	db := infrastructure.NewDB(infrastructure.DBConfig{
		Host:     "localhost",
		Port:     "15432",
		DBName:   "echo_example",
		User:     "root",
		Password: "password",
	})
	router := api.NewRouter(db)

	router.Logger.Fatal(router.Start(":1323"))
}
