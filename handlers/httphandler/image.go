package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/geekakili/portside/driver"
	"github.com/geekakili/portside/models"
	repository "github.com/geekakili/portside/repository/image"
	"github.com/go-chi/chi"
	validate "gopkg.in/dealancer/validate.v2"
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
	Image  string   `validate:"empty=false"`
	Labels []string `validate:"empty=false"`
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
	labels, ok := r.URL.Query()["withLabel"]
	if ok && len(labels) > 0 {
		var dockerImages []types.ImageInspect
		for _, label := range labels {
			labeledImages, err := image.repo.GetByLabel(r.Context(), label)
			if err == nil && len(labeledImages) > 0 {
				for _, labeledImage := range labeledImages {
					dockerImage, _, err := image.dockerClient.ImageInspectWithRaw(r.Context(), labeledImage)
					if err != nil {
						continue
					}
					dockerImages = append(dockerImages, dockerImage)
				}
			}
		}
		if len(dockerImages) > 0 {
			image.marshallDockerImages(r.Context(), w, dockerImages)
		} else {
			respondWithJSON(w, http.StatusNotFound, "No images with matching labels found on the host system")
		}
	} else {
		images, err := image.dockerClient.ImageList(r.Context(), types.ImageListOptions{})
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, "Opps, Something went wrong")
		}
		if len(images) > 0 {
			//dockerImages = append(dockerImages, images)
			image.marshallDockerImages(r.Context(), w, images)
		} else {
			respondWithJSON(w, http.StatusNotFound, "No images found on the host system")
		}
	}
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
		return
	}

	err = validate.Validate(imageData)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Image name is missing, check your request and try again")
		return
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
		return
	}

	err = validate.Validate(labelInfo)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Some data is missing, check your request and try again")
		return
	}
	err = image.repo.AddLabel(r.Context(), labelInfo.Image, labelInfo.Labels...)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	labels, err := image.repo.GetImageLabels(r.Context(), labelInfo.Image)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, labels)
}

func (image *imageHandler) marshallDockerImages(ctx context.Context, w http.ResponseWriter, dockerImages interface{}) {
	var images []models.Image
	dockerImagesVal := reflect.ValueOf(dockerImages)
	for i := 0; i < dockerImagesVal.Len(); i++ {
		var name string
		var tag string
		imageData := dockerImagesVal.Index(i)
		var dockerImageName string
		repoTags := imageData.FieldByName("RepoTags").Interface().([]string)
		if len(repoTags) > 0 {
			dockerImageName = repoTags[0]
			repoData := strings.Split(dockerImageName, ":")
			name = repoData[0]
			tag = repoData[1]
		}

		imageID := imageData.FieldByName("ID").String()
		imageSize := imageData.FieldByName("Size").Int()
		repoDigests := imageData.FieldByName("RepoDigests").Interface().([]string)
		dockerImage := models.Image{
			Name:       name,
			ID:         imageID,
			Size:       imageSize,
			Repository: "repo",
			Tag:        tag,
			Digests:    repoDigests,
			Labels:     make([]string, 0),
		}

		if len(dockerImageName) > 0 {
			labels, err := image.repo.GetImageLabels(ctx, dockerImageName)
			if err == nil {
				dockerImage.Labels = labels
			}
		}
		images = append(images, dockerImage)
	}
	respondWithJSON(w, http.StatusOK, images)
}
