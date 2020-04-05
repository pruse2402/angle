package handlers

import (
	"angle/src/datastore"
	"angle/src/errs"
	"angle/src/internal"
	"angle/src/models"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func (p *Provider) InsertRawMaterial(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	rawMaterialIns := models.RawMaterial{}

	if !parseJson(rw, r.Body, &rawMaterialIns) {
		return
	}

	hasErr, validationErr := internal.RawValidation(p.db, rawMaterialIns)
	if hasErr {
		log.Printf("ERROR: InsertRawMaterial - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	rawMaterialIns.Id = bson.NewObjectId()
	rawMaterialIns.DateCreated = time.Now().UTC()
	rawMaterialIns.LastUpdated = time.Now().UTC()

	if !isObjectIDValid(rawMaterialIns.Vendors.VendorCode.Hex()) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid vendor ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid vendor ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	var resp string
	rawMaterial, err := datastore.NewRawMaterial(dbSession).InsertRawMaterial(rawMaterialIns)
	if err != nil {
		resp = "Error while saving raw material"
		renderJson(rw, http.StatusUnauthorized, resp)
		return
	}

	res := struct {
		ID      bson.ObjectId `json:"id"`
		Message string        `json:"message"`
	}{
		rawMaterial.Id,
		`Raw Material saved successfully `,
	}

	renderJson(rw, http.StatusOK, res)
	return

}

func (p *Provider) GetRawMaterials(rw http.ResponseWriter, r *http.Request) {
	dbSession := p.db.Copy()
	defer dbSession.Close()

	rawMaterials, err := datastore.NewRawMaterial(dbSession).FindAllRawMaterial()
	if err != nil {
		log.Printf("ERROR: GetRawMaterialDetailList %s", err)
		return
	}

	resp := struct {
		RawMaterials *[]models.RawMaterial `json:"rawMaterials"`
	}{
		rawMaterials,
	}

	renderJson(rw, http.StatusOK, resp)

}
