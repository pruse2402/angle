package internal

import (
	"angle/src/datastore"
	"angle/src/models"

	"gopkg.in/mgo.v2"
)

func EmployeeValidation(dbSession *mgo.Session, employeeDetails models.EmployeeDetails) (bool, map[string]interface{}) {

	employeeDetailsIns := datastore.NewEmployeeDetails(dbSession)

	hasErr, validationErr := employeeDetailsIns.Validate(employeeDetails)

	return hasErr, validationErr
}
