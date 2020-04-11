package handlers

import (
	"angle/src/datastore"
	"angle/src/errs"
	"angle/src/internal"
	"angle/src/models"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/mgo.v2/bson"
)

func (p *Provider) InsertVendorDetails(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	vendorDetailsIns := models.VendorDetails{}

	if !parseJson(rw, r.Body, &vendorDetailsIns) {
		return
	}

	hasErr, validationErr := internal.Validation(p.db, vendorDetailsIns)
	if hasErr {
		log.Printf("ERROR: InsertVendorDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	vendorDetailsIns.Id = bson.NewObjectId()
	vendorDetailsIns.DateCreated = time.Now().UTC()
	vendorDetailsIns.LastUpdated = time.Now().UTC()

	// c := datastore.NewVendorDetails(dbSession).FindCount(vendorDetailsIns)
	// vendorDetailsIns.SequenceNumber = c

	var resp string
	vendorDetails, err := datastore.NewVendorDetails(dbSession).InsertVendorDetails(vendorDetailsIns)
	if err != nil {
		resp = "Error while saving vendor detail"
		renderJson(rw, http.StatusUnauthorized, resp)
		return
	}

	res := struct {
		ID      bson.ObjectId `json:"id"`
		Message string        `json:"message"`
	}{
		vendorDetails.Id,
		`vendor detail saved successfully `,
	}

	renderJson(rw, http.StatusOK, res)
	return

}

func (p *Provider) GetVendors(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	vendorDetails, err := datastore.NewVendorDetails(dbSession).FindAllVendor()
	if err != nil {
		log.Printf("ERROR: GetVendorDetailList %s", err)
		return
	}

	sort.Sort(models.VendorDetailsByName(*vendorDetails))

	resp := struct {
		VendorDetails *[]models.VendorDetails `json:"vendorDetails"`
	}{
		vendorDetails,
	}

	renderJson(rw, http.StatusOK, resp)

}

func (p *Provider) UpdateVendorDetails(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	vendorID := chi.URLParam(r, "id")

	if !isObjectIDValid(vendorID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid vendor ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid vendor ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	vendorDetails, err := datastore.NewVendorDetails(dbSession).FindByID(bson.ObjectIdHex(vendorID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find vendor details!.",
		}
		renderError(rw, err)
		return
	}

	if !parseJson(rw, r.Body, &vendorDetails) {
		return
	}

	hasErr, validationErr := internal.Validation(p.db, *vendorDetails)
	if hasErr {
		log.Printf("ERROR: InsertVendorDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	vendorDetails.LastUpdated = time.Now().UTC()

	if err = datastore.NewVendorDetails(dbSession).Update(vendorDetails.Id, *vendorDetails); err != nil {
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
		vendorDetails.Id,
		"Vendor Details updated successfully.",
	}

	renderJson(rw, http.StatusOK, res)

}

func (p *Provider) RemoveVendorDetails(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	vendorID := chi.URLParam(r, "id")
	if !isObjectIDValid(vendorID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid vendor ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid vendor ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	vendorDetails, err := datastore.NewVendorDetails(dbSession).FindByID(bson.ObjectIdHex(vendorID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find vendor details!.",
		}
		renderError(rw, err)
		return
	}

	err = internal.DeleteVendorDetails(p.db, vendorID, "vendorDetails", vendorDetails)
	if err != nil {
		p.log.Printf("ERROR: Handler - Delete - %q\n", err)
		renderError(rw, err)
		return
	}

	// Constructing response for client
	res := struct {
		Message string `json:"message"`
	}{
		"Vendor Details Removed Successfully.",
	}

	renderJson(rw, http.StatusOK, res)
}
