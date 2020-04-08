package internal

import (
	"angle/src/datastore"
	"angle/src/models"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func RawValidation(dbSession *mgo.Session, rawMaterial models.RawMaterial) (bool, map[string]interface{}) {

	rawMaterialIns := datastore.NewRawMaterial(dbSession)

	hasErr, validationErr := rawMaterialIns.Validate(rawMaterial)

	return hasErr, validationErr
}

func DeleteRawMaterial(dbSession *mgo.Session, rawId, collection string, rawMaterial *models.RawMaterial) error {

	backUp := struct {
		Id         bson.ObjectId      `bson:"_id"`
		Module     string             `bson:"module"`
		BackupDate time.Time          `bson:"backupDate"`
		Ids        string             `bson:"ids"`
		Data       models.RawMaterial `bson:"data"`
	}{
		bson.NewObjectId(),
		collection,
		time.Now(),
		rawId,
		*rawMaterial,
	}
	if err := dbSession.DB("angle").C("BackUp").Insert(&backUp); err != nil {
		log.Printf("ERROR: Raw Material Backup (%s) - %s", backUp.Module, err)
		return err
	}
	if err := dbSession.DB("angle").C(collection).Remove(bson.M{"_id": bson.ObjectIdHex(rawId)}); err != nil {
		log.Printf("ERROR: Raw Material Details Remove (%s) - %s", backUp.Module, err)
		return err
	}
	return nil
}
