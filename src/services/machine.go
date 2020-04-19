package services

import (
	"angle/src/models"

	"gopkg.in/mgo.v2/bson"
)

type MachineDetails interface {
	Validate(models.MachineDetails) (bool, map[string]interface{})
	//FindCount(models.VendorDetails) int
	InsertMachineDetails(models.MachineDetails) (*models.MachineDetails, error)
	FindAllMachineDetails() (*[]models.MachineDetails, error)
	FindByID(bson.ObjectId) (*models.MachineDetails, error)
	Update(bson.ObjectId, models.MachineDetails) error
}
