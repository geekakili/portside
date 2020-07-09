package repository

import (
	"context"

	"github.com/geekakili/portside/models"
	"github.com/zippoxer/bow"
)

// NewBadgerLabelRepo initializes Label repo
func NewBadgerLabelRepo(badgerDBConn *bow.DB) LabelRepository {
	return &badgerDB{Conn: badgerDBConn}
}

type badgerDB struct {
	Conn *bow.DB
}

func (db *badgerDB) AddLabel(ctx context.Context, label string, description string) (*models.Label, error) {
	labelData := new(models.Label)
	err := db.Conn.Bucket("labels").Get(label, labelData)
	if err != nil {
		labelData.Name = label
		labelData.Description = description
		labelData.Images = make([]string, 0)
		err = db.Conn.Bucket("labels").Put(labelData)
		if err != nil {
			return nil, err
		}
	}

	return labelData, nil
}

func (db *badgerDB) GetLabel(ctx context.Context, label string) (models.Label, error) {
	labelData := new(models.Label)
	err := db.Conn.Bucket("labels").Get(label, labelData)
	return *labelData, err
}

func (db *badgerDB) GetLabels(ctx context.Context) []models.Label {
	labelData := new(models.Label)
	var labels []models.Label
	iter := db.Conn.Bucket("labels").Iter()
	defer iter.Close()
	for iter.Next(labelData) {
		labels = append(labels, *labelData)
	}
	return labels
}
