package repository

import (
	"context"

	"github.com/geekakili/portside/models"
	"github.com/zippoxer/bow"
)

// NewBadgerImageRepo initailizes image repository
func NewBadgerImageRepo(badgerDBConn *bow.DB) ImageRepository {
	return &badgerDB{Conn: badgerDBConn}
}

type badgerDB struct {
	Conn *bow.DB
}

type labelData struct {
	Id string `bow:"key"`
}

func (db *badgerDB) AddLabel(ctx context.Context, tag string, labels ...string) error {
	label := models.ImageLabel{
		Id:     tag,
		Labels: labels,
	}

	for _, labelName := range labels {
		label := labelData{
			Id: labelName,
		}
		err := db.Conn.Bucket("labels").Put(label)
		if err != nil {
			return err
		}
	}

	err := db.Conn.Bucket("labeledImages").Put(label)
	return err
}

// GetImageLabels Returns a list of labels associated with the image
func (db *badgerDB) GetImageLabels(ctx context.Context, imageName string) (labels []string, err error) {
	var imageLabel models.ImageLabel
	err = db.Conn.Bucket("labeledImages").Get(imageName, &imageLabel)
	if err != nil {
		return nil, err
	}
	return imageLabel.Labels, nil
}

func (db *badgerDB) GetByLabel(ctx context.Context, label string) ([]*models.Image, error) {
	return nil, nil
}
