package main

import "github.com/ricky2122/go-echo-example/infrastructure/api"

func main() {
	router := api.NewRouter()

	router.Logger.Fatal(router.Start(":1323"))
}
