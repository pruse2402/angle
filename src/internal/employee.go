package internal

import (
	"angle/src/datastore"
	"angle/src/models"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func EmployeeValidation(dbSession *mgo.Session, employeeDetails models.EmployeeDetails) (bool, map[string]interface{}) {

	employeeDetailsIns := datastore.NewEmployeeDetails(dbSession)

	hasErr, validationErr := employeeDetailsIns.Validate(employeeDetails)

	return hasErr, validationErr
}

func DeleteEmployee(dbSession *mgo.Session, empId, collection string, empDetails *models.EmployeeDetails) error {

	backUp := struct {
		Id         bson.ObjectId          `bson:"_id"`
		Module     string                 `bson:"module"`
		BackupDate time.Time              `bson:"backupDate"`
		Ids        string                 `bson:"ids"`
		Data       models.EmployeeDetails `bson:"data"`
	}{
		bson.NewObjectId(),
		collection,
		time.Now(),
		empId,
		*empDetails,
	}
	if err := dbSession.DB("angle").C("BackUp").Insert(&backUp); err != nil {
		log.Printf("ERROR: Employee Details Backup (%s) - %s", backUp.Module, err)
		return err
	}
	if err := dbSession.DB("angle").C(collection).Remove(bson.M{"_id": bson.ObjectIdHex(empId)}); err != nil {
		log.Printf("ERROR: Employee Details Remove (%s) - %s", backUp.Module, err)
		return err
	}
	return nil
}
