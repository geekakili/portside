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

// Addlabel adds a label to the database
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

// GetLabel returns a label from the database
func (db *badgerDB) GetLabel(ctx context.Context, label string) (models.Label, error) {
	labelData := new(models.Label)
	err := db.Conn.Bucket("labels").Get(label, labelData)
	return *labelData, err
}

// GetLabels returns labels from the database
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

// Delete deletes a label from the database
func (db *badgerDB) Delete(ctx context.Context, label string) (bool, error) {
	labelData, err := db.GetLabel(ctx, label)
	if err == nil {
		err := db.Conn.Bucket("labels").Delete(labelData.Name)
		if err != nil {
			return false, err
		}
		return true, err
	}
	return false, err
}
