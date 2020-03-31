package services

import (
	"angle/src/models"

	"gopkg.in/mgo.v2/bson"
)

type VendorDetails interface {
	Validate(models.VendorDetails) (bool, map[string]interface{})
	//FindCount(models.VendorDetails) int
	InsertVendorDetails(models.VendorDetails) (*models.VendorDetails, error)
	FindAllVendor() (*[]models.VendorDetails, error)
	FindByID(bson.ObjectId) (*models.VendorDetails, error)
	Update(bson.ObjectId, models.VendorDetails) error
}
