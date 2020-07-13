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
	router.Get("/list/{label}", handler.list)
	router.Get("/list", handler.list)
	router.Put("/update/{label}", handler.update)
	router.Delete("/delete/{label}", handler.delete)
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

// list ...
func (label *labelHandler) list(w http.ResponseWriter, r *http.Request) {
	labelName := chi.URLParam(r, "label")
	if len(labelName) > 0 {
		labelData, err := label.repo.GetLabel(r.Context(), labelName)
		if err != nil {
			respondWithJSON(w, http.StatusNotFound, "Label with such a name not found")
			return
		}
		respondWithJSON(w, http.StatusOK, labelData)
	} else {
		labels := label.repo.GetLabels(r.Context())
		if len(labels) > 0 {
			respondWithJSON(w, http.StatusOK, labels)
			return
		}
		respondWithJSON(w, http.StatusNotFound, "No labels exist on this host")
	}
}

func (label *labelHandler) update(w http.ResponseWriter, r *http.Request) {
	labelName := chi.URLParam(r, "label")
	labelData := new(models.Label)
	err := json.NewDecoder(r.Body).Decode(labelData)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Error processing label data")
		return
	}
	if len(labelName) > 0 {
		if len(labelData.Name) > 0 || len(labelData.Description) > 0 {
			err := label.repo.Updatelabel(r.Context(), labelName, *labelData)
			if err != nil {
				respondWithJSON(w, http.StatusInternalServerError, "Could not update label")
				return
			}
			respondWithJSON(w, http.StatusOK, "Label updated successfully")
			return
		}
		respondWithJSON(w, http.StatusBadRequest, "Could not update label, no label data is provided")
		return
	}
	respondWithJSON(w, http.StatusBadRequest, "Could not update label, label not found")
}

func (label *labelHandler) delete(w http.ResponseWriter, r *http.Request) {
	labelName := chi.URLParam(r, "label")
	if len(labelName) > 0 {
		deleted, err := label.repo.Delete(r.Context(), labelName)
		if err != nil || !deleted {
			respondWithJSON(w, http.StatusInternalServerError, "Label failed to delete")
			return
		}
		respondWithJSON(w, http.StatusOK, "Label deleted successfully")
		return
	}
	respondWithJSON(w, http.StatusBadRequest, "Could not delete label, request is malformed")
}
