package services

import (
	"angle/src/models"

	"gopkg.in/mgo.v2/bson"
)

type CustomerDetails interface {
	Validate(models.CustomerDetails) (bool, map[string]interface{})
	//FindCount(models.VendorDetails) int
	InsertCustomerDetails(models.CustomerDetails) (*models.CustomerDetails, error)
	FindAllCustomer() (*[]models.CustomerDetails, error)
	FindByID(bson.ObjectId) (*models.CustomerDetails, error)
	Update(bson.ObjectId, models.CustomerDetails) error
}
