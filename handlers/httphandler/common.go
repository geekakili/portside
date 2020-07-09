package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type httpResponse struct {
	Status  int
	Message interface{}
}

// respondWithJSON returns a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(httpResponse{Status: code, Message: payload})
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func setupRoute(parent *chi.Mux, child *chi.Mux, prefix string) {
	parent.Route(prefix, func(rt chi.Router) {
		rt.Mount("/", child)
	})
}
