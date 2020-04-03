package internal

import (
	"angle/src/datastore"
	"angle/src/models"

	"gopkg.in/mgo.v2"
)

func RawValidation(dbSession *mgo.Session, rawMaterial models.RawMaterial) (bool, map[string]interface{}) {

	rawMaterialIns := datastore.NewRawMaterial(dbSession)

	hasErr, validationErr := rawMaterialIns.Validate(rawMaterial)

	return hasErr, validationErr
}
