package services

import (
	"angle/src/models"

	"gopkg.in/mgo.v2/bson"
)

type EmployeeDetails interface {
	Validate(models.EmployeeDetails) (bool, map[string]interface{})
	//FindCount(models.VendorDetails) int
	InsertEmployeeDetails(models.EmployeeDetails) (*models.EmployeeDetails, error)
	FindAllEmployee() (*[]models.EmployeeDetails, error)
	FindByID(bson.ObjectId) (*models.EmployeeDetails, error)
	Update(bson.ObjectId, models.EmployeeDetails) error
}
