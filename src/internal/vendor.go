package internal

import (
	"angle/src/datastore"
	"angle/src/models"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Validation(dbSession *mgo.Session, vendorDetails models.VendorDetails) (bool, map[string]interface{}) {

	vendorDetailsIns := datastore.NewVendorDetails(dbSession)

	hasErr, validationErr := vendorDetailsIns.Validate(vendorDetails)

	return hasErr, validationErr
}

func DeleteVendorDetails(dbSession *mgo.Session, vendorId, collection string, vendorDetails *models.VendorDetails) error {

	backUp := struct {
		Id         bson.ObjectId        `bson:"_id"`
		Module     string               `bson:"module"`
		BackupDate time.Time            `bson:"backupDate"`
		Vendor     string               `bson:"vendorId"`
		Data       models.VendorDetails `bson:"data"`
	}{
		bson.NewObjectId(),
		collection,
		time.Now(),
		vendorId,
		*vendorDetails,
	}
	if err := dbSession.DB("angle").C("BackUp").Insert(&backUp); err != nil {
		log.Printf("ERROR: Vendor Backup (%s) - %s", backUp.Module, err)
		return err
	}
	if err := dbSession.DB("angle").C(collection).Remove(bson.M{"_id": bson.ObjectIdHex(vendorId)}); err != nil {
		log.Printf("ERROR: Vendor Details Remove (%s) - %s", backUp.Module, err)
		return err
	}
	return nil
}
