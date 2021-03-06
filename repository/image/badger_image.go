package repository

import (
	"context"

	"github.com/geekakili/portside/models"
	"github.com/geekakili/portside/repository"
	"github.com/zippoxer/bow"
)

// NewBadgerImageRepo initailizes image repository
func NewBadgerImageRepo(badgerDBConn *bow.DB) ImageRepository {
	return &badgerDB{Conn: badgerDBConn}
}

type badgerDB struct {
	Conn *bow.DB
}

func (db *badgerDB) AddLabel(ctx context.Context, tag string, labels ...string) error {
	imageLabels, _ := db.GetImageLabels(ctx, tag)
	label := models.ImageLabel{
		Image:  tag,
		Labels: imageLabels,
	}

	for _, labelName := range labels {
		var newLabel models.Label
		err := db.Conn.Bucket("labels").Get(labelName, &newLabel)
		if err != nil {
			newLabel = models.Label{
				Name: labelName,
			}
		}

		imageFound := repository.ArrayContains(newLabel.Images, tag)
		if !imageFound {
			newLabel.Images = append(newLabel.Images, tag)
			err := db.Conn.Bucket("labels").Put(newLabel)
			if err != nil {
				return err
			}
		}

		labelFound := repository.ArrayContains(imageLabels, labelName)
		if !labelFound {
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

func (db *badgerDB) GetByLabel(ctx context.Context, label string) ([]string, error) {
	var imageLabel models.Label
	err := db.Conn.Bucket("labels").Get(label, &imageLabel)
	return imageLabel.Images, err
}
