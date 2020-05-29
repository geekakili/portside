package httphandler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/geekakili/portside/driver"
	"github.com/geekakili/portside/models"
	repository "github.com/geekakili/portside/repository/image"
	"github.com/go-chi/chi"
)

// SetupImageHandler setups routes to handle image requests
func setupImageHandler(db *driver.DB, client *client.Client, httpRouter *chi.Mux) {
	handler := &imageHandler{
		repo:         repository.NewBadgerImageRepo(db.Badger),
		dockerClient: client,
	}

	router := chi.NewRouter()
	router.Get("/list", handler.list)
	setupRoute(httpRouter, router, "/image")
}

// imageHandler ...
type imageHandler struct {
	repo         repository.ImageRepository
	dockerClient *client.Client
}

// list returns a list of all docker images on the host machine
func (image *imageHandler) list(w http.ResponseWriter, r *http.Request) {
	dockerImages, err := image.dockerClient.ImageList(r.Context(), types.ImageListOptions{})
	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, http.StatusInternalServerError, "Opps, Something went wrong")
	}

	var images []models.Image
	if len(dockerImages) > 0 {
		for _, imageData := range dockerImages {
			repoData := strings.Split(imageData.RepoTags[0], ":")
			repo := repoData[0]
			tag := repoData[1]
			dockerImage := models.Image{
				ID:         imageData.ID,
				Size:       imageData.Size,
				Repository: repo,
				Tag:        tag,
				Digests:    imageData.RepoDigests,
			}
			images = append(images, dockerImage)
		}
	}
	respondWithJSON(w, http.StatusOK, images)
}

// GetByName ...
func (image *imageHandler) getByName(ctx context.Context, name string) {

}

// GetByLabel ...
func (image *imageHandler) getByLabel(ctx context.Context, label string) {

}

// PullImage ...
func (image *imageHandler) pullImage(ctx context.Context, name string) {

}
