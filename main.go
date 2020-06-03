package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/geekakili/portside/driver"
	"github.com/geekakili/portside/handlers/httphandler"
)

func main() {
	dbConnection, err := driver.ConnectBadger("./badger")
	defer dbConnection.Badger.Close()

	if err != nil {
		log.Fatal(err)
	}

	router, err := httphandler.SetupRoutes(dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Portside server listening on port 8005")
	err = http.ListenAndServe(":8005", router)
	if err != nil {
		log.Fatal(err)
	}

}
