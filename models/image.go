package models

// Image holds basic information about the docker image
type Image struct {
	ID         string   //ID of the docker image
	Size       int64    // Size of the docker image
	Repository string   // Repository of the docker image
	Tag        string   //Tag of the docker image
	Digests    []string // List of sha256 digests
}
