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
	Id     string `bow:"key"`
	Images []string
}

func (db *badgerDB) AddLabel(ctx context.Context, tag string, labels ...string) error {
	imageLabels, _ := db.GetImageLabels(ctx, tag)
	label := models.ImageLabel{
		Id:     tag,
		Labels: imageLabels,
	}

	for _, labelName := range labels {
		var newLabel labelData
		err := db.Conn.Bucket("labels").Get(labelName, &newLabel)
		if err != nil {
			newLabel = labelData{
				Id: labelName,
			}
		}

		imageFound := ArrayContains(newLabel.Images, tag)
		if imageFound == false {
			newLabel.Images = append(newLabel.Images, tag)
			err := db.Conn.Bucket("labels").Put(newLabel)
			if err != nil {
				return err
			}
		}

		labelFound := ArrayContains(imageLabels, labelName)
		if labelFound == false {
			label.Labels = append(label.Labels, labelName)
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
