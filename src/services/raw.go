package services

import (
	"angle/src/models"

	"gopkg.in/mgo.v2/bson"
)

type RawMaterial interface {
	Validate(models.RawMaterial) (bool, map[string]interface{})
	InsertRawMaterial(models.RawMaterial) (*models.RawMaterial, error)
	FindAllRawMaterial() (*[]models.RawMaterial, error)
	FindByID(bson.ObjectId) (*models.RawMaterial, error)
	Update(bson.ObjectId, models.RawMaterial) error
}
