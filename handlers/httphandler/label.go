package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/geekakili/portside/driver"
	"github.com/geekakili/portside/models"
	repository "github.com/geekakili/portside/repository/label"
	"github.com/go-chi/chi"
	"gopkg.in/dealancer/validate.v2"
)

// setupLabelHandler setups routes to handle label requests
func setupLabelHandler(db *driver.DB, client *client.Client, httpRouter *chi.Mux) {
	handler := &labelHandler{
		repo: repository.NewBadgerLabelRepo(db.Badger),
	}

	router := chi.NewRouter()
	router.Post("/add", handler.add)
	// router.Post("/list", handler.list)
	// router.Post("/edit", handler.edit)
	// router.Post("/delete", handler.delete)
	setupRoute(httpRouter, router, "/labels")
}

// labelHandler ...
type labelHandler struct {
	repo repository.LabelRepository
}

// list returns a list of all docker images on the host machine
func (label *labelHandler) add(w http.ResponseWriter, r *http.Request) {
	labelData := new(models.Label)
	err := json.NewDecoder(r.Body).Decode(labelData)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, "Couldn't parse label data")
	}

	err = validate.Validate(labelData)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Label name is missing")
		return
	}

	newlabel, err := label.repo.AddLabel(r.Context(), labelData.Name, labelData.Description)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, "An error occured while processing label")
		return
	}
	respondWithJSON(w, http.StatusOK, newlabel)
}

// // GetByName ...
// func (label *labelHandler) list(w http.ResponseWriter, r *http.Request) {

// }

// // GetByLabel ...
// func (label *labelHandler) edit(w http.ResponseWriter, r *http.Request) {

// }

// // PullImage pulls image from remote repository
// func (label *labelHandler) delete(w http.ResponseWriter, r *http.Request) {

// }
