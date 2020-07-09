package httphandler

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/geekakili/portside/driver"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// SetupRoutes sets up all routes
func SetupRoutes(db *driver.DB) (*chi.Mux, error) {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	client, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	setupImageHandler(db, client, router)
	setupLabelHandler(db, client, router)

	return router, nil
}
