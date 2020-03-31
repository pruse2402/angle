package internal

import (
	"angle/src/datastore"
	"angle/src/models"

	"gopkg.in/mgo.v2"
)

func Validation(dbSession *mgo.Session, vendorDetails models.VendorDetails) (bool, map[string]interface{}) {

	vendorDetailsIns := datastore.NewVendorDetails(dbSession)

	hasErr, validationErr := vendorDetailsIns.Validate(vendorDetails)

	return hasErr, validationErr
}
