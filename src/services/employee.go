package services

import "angle/src/models"

type EmployeeDetails interface {
	Validate(models.EmployeeDetails) (bool, map[string]interface{})
	//FindCount(models.VendorDetails) int
	InsertEmployeeDetails(models.EmployeeDetails) (*models.EmployeeDetails, error)
	FindAllEmployee() (*[]models.EmployeeDetails, error)
	// FindByID(bson.ObjectId) (*models.CustomerDetails, error)
	// Update(bson.ObjectId, models.CustomerDetails) error
}
