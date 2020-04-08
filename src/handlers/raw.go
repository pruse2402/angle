package handlers

import (
	"angle/src/datastore"
	"angle/src/errs"
	"angle/src/internal"
	"angle/src/models"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
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

func (p *Provider) UpdateRawMaterial(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	rawID := chi.URLParam(r, "id")

	if !isObjectIDValid(rawID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid raw material ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid raw material ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	rawMaterials, err := datastore.NewRawMaterial(dbSession).FindByID(bson.ObjectIdHex(rawID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find raw material details!.",
		}
		renderError(rw, err)
		return
	}

	if !parseJson(rw, r.Body, &rawMaterials) {
		return
	}

	hasErr, validationErr := internal.RawValidation(p.db, *rawMaterials)
	if hasErr {
		log.Printf("ERROR: UpdateRawMaterialDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	rawMaterials.LastUpdated = time.Now().UTC()

	if err = datastore.NewRawMaterial(dbSession).Update(rawMaterials.Id, *rawMaterials); err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong!. Please try again!.",
		}
		renderError(rw, err)
		return
	}

	// Constructing response for client
	res := struct {
		ID      bson.ObjectId `json:"id"`
		Message string        `json:"message"`
	}{
		rawMaterials.Id,
		"raw materials updated successfully.",
	}

	renderJson(rw, http.StatusOK, res)
}

func (p *Provider) RemoveRawMaterial(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	rawID := chi.URLParam(r, "id")
	if !isObjectIDValid(rawID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid Raw ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid Raw ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	rawMaterial, err := datastore.NewRawMaterial(dbSession).FindByID(bson.ObjectIdHex(rawID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find raw material details!.",
		}
		renderError(rw, err)
		return
	}

	err = internal.DeleteRawMaterial(p.db, rawID, "rawMaterial", rawMaterial)
	if err != nil {
		p.log.Printf("ERROR: Handler - Delete - %q\n", err)
		renderError(rw, err)
		return
	}

	// Constructing response for client
	res := struct {
		Message string `json:"message"`
	}{
		"Raw Material Details Removed Successfully.",
	}

	renderJson(rw, http.StatusOK, res)
}
