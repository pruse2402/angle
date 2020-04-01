package internal

import (
	"angle/src/datastore"
	"angle/src/models"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CustomerValidation(dbSession *mgo.Session, customerDetails models.CustomerDetails) (bool, map[string]interface{}) {

	customerDetailsIns := datastore.NewCustomerDetails(dbSession)

	hasErr, validationErr := customerDetailsIns.Validate(customerDetails)

	return hasErr, validationErr
}

func DeleteCustomerDetails(dbSession *mgo.Session, customerId, collection string, customerDetails *models.CustomerDetails) error {

	backUp := struct {
		Id         bson.ObjectId          `bson:"_id"`
		Module     string                 `bson:"module"`
		BackupDate time.Time              `bson:"backupDate"`
		Ids        string                 `bson:"ids"`
		Data       models.CustomerDetails `bson:"data"`
	}{
		bson.NewObjectId(),
		collection,
		time.Now(),
		customerId,
		*customerDetails,
	}
	if err := dbSession.DB("angle").C("BackUp").Insert(&backUp); err != nil {
		log.Printf("ERROR: Customer Backup (%s) - %s", backUp.Module, err)
		return err
	}
	if err := dbSession.DB("angle").C(collection).Remove(bson.M{"_id": bson.ObjectIdHex(customerId)}); err != nil {
		log.Printf("ERROR: Customer Details Remove (%s) - %s", backUp.Module, err)
		return err
	}
	return nil
}
