package services

import "angle/src/models"

type RawMaterial interface {
	Validate(models.RawMaterial) (bool, map[string]interface{})
	InsertRawMaterial(models.RawMaterial) (*models.RawMaterial, error)
}
