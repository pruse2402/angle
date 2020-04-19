package internal

import (
	"angle/src/datastore"
	"angle/src/models"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func MachineValidation(dbSession *mgo.Session, machineDetails models.MachineDetails) (bool, map[string]interface{}) {

	machineDetailsIns := datastore.NewMachineDetails(dbSession)

	hasErr, validationErr := machineDetailsIns.Validate(machineDetails)

	return hasErr, validationErr
}

func DeleteMachine(dbSession *mgo.Session, machineId, collection string, machineDetails *models.MachineDetails) error {

	backUp := struct {
		Id         bson.ObjectId         `bson:"_id"`
		Module     string                `bson:"module"`
		BackupDate time.Time             `bson:"backupDate"`
		Ids        string                `bson:"ids"`
		Data       models.MachineDetails `bson:"data"`
	}{
		bson.NewObjectId(),
		collection,
		time.Now(),
		machineId,
		*machineDetails,
	}
	if err := dbSession.DB("angle").C("BackUp").Insert(&backUp); err != nil {
		log.Printf("ERROR: Machine Details Backup (%s) - %s", backUp.Module, err)
		return err
	}
	if err := dbSession.DB("angle").C(collection).Remove(bson.M{"_id": bson.ObjectIdHex(machineId)}); err != nil {
		log.Printf("ERROR: Machine Details Remove (%s) - %s", backUp.Module, err)
		return err
	}
	return nil
}
