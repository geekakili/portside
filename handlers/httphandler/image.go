package httphandler

import (
	"context"
	"encoding/json"
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

// ImageRepository ..
type ImageRepository string

// Dockerhub ...
const (
	Dockerhub ImageRepository = "docker.io/library"
)

type progress struct {
	Status         string
	ProgressDetail struct {
		current int
		Total   int
	}
	Progress string
	Digest   string
	ID       string
}

type imageLabel struct {
	Image  string
	Labels []string
}

// SetupImageHandler setups routes to handle image requests
func setupImageHandler(db *driver.DB, client *client.Client, httpRouter *chi.Mux) {
	handler := &imageHandler{
		repo:         repository.NewBadgerImageRepo(db.Badger),
		dockerClient: client,
	}

	router := chi.NewRouter()
	router.Get("/list", handler.list)
	router.Post("/pull", handler.pullImage)
	router.Post("/label", handler.labelImage)
	setupRoute(httpRouter, router, "/images")
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
			fmt.Println(imageData.ID)
			images = append(images, dockerImage)
		}
	}
	respondWithJSON(w, http.StatusOK, images)
}

// GetByName ...
func (image *imageHandler) getByName(w http.ResponseWriter, r *http.Request) {

}

// GetByLabel ...
func (image *imageHandler) getByLabel(ctx context.Context, label string) {

}

// PullImage pulls image from remote repository
func (image *imageHandler) pullImage(w http.ResponseWriter, r *http.Request) {
	imageData := new(models.Image)
	err := json.NewDecoder(r.Body).Decode(imageData)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, "Couldn't parse image name")
	}

	var remoteImage string
	if len(imageData.Repository) == 0 {
		remoteImage = fmt.Sprintf("%s/%s", Dockerhub, imageData.Name)
	} else {
		remoteImage = fmt.Sprintf("%s/%s", imageData.Repository, imageData.Name)
	}

	if len(imageData.Tag) > 0 {
		remoteImage = fmt.Sprintf("%s:%s", remoteImage, imageData.Tag)
	}

	fmt.Println(remoteImage)
	reader, err := image.dockerClient.ImagePull(r.Context(), remoteImage, types.ImagePullOptions{})
	if err != nil {
		errString := err.Error()
		respondWithJSON(w, http.StatusInternalServerError, errString)
		return
	}
	defer reader.Close()

	buff := make([]byte, 1024)
	lastResponse := new(progress)
	for {
		_, err := reader.Read(buff)
		if err != nil {
			break
		}
		status := strings.Split(string(buff), "\n")
		json.Unmarshal([]byte(status[0]), lastResponse)
	}

	inspectedImage, _, err := image.dockerClient.ImageInspectWithRaw(r.Context(), imageData.Name)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, inspectedImage)
}

// labelImage ...
func (image *imageHandler) labelImage(w http.ResponseWriter, r *http.Request) {
	labelInfo := new(imageLabel)
	err := json.NewDecoder(r.Body).Decode(labelInfo)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, "Oops, something went wrong")
	}
	err = image.repo.AddLabel(r.Context(), labelInfo.Image, labelInfo.Labels...)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
	}
	respondWithJSON(w, http.StatusOK, labelInfo)
}
