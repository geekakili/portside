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
