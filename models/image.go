package models

// Image holds basic information about the docker image
type Image struct {
	Name       string   `validate:"empty=false"` // Name of image
	ID         string   // ID of the docker image
	Size       int64    // Size of the docker image
	Repository string   // Repository of the docker image
	Tag        string   // Tag of the docker image
	Digests    []string // List of sha256 digests
	Labels     []string // List of image labels
}

// ImageLabel holds information needed to label an image
type ImageLabel struct {
	Id     string `bow:"key"`
	Labels []string
}
