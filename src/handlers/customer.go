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

func (p *Provider) InsertCustomerDetails(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	customerDetailsIns := models.CustomerDetails{}

	if !parseJson(rw, r.Body, &customerDetailsIns) {
		return
	}

	hasErr, validationErr := internal.CustomerValidation(p.db, customerDetailsIns)
	if hasErr {
		log.Printf("ERROR: InsertCustomerDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	customerDetailsIns.Id = bson.NewObjectId()
	customerDetailsIns.DateCreated = time.Now().UTC()
	customerDetailsIns.LastUpdated = time.Now().UTC()

	// c := datastore.NewVendorDetails(dbSession).FindCount(vendorDetailsIns)
	// vendorDetailsIns.SequenceNumber = c

	var resp string
	customerDetails, err := datastore.NewCustomerDetails(dbSession).InsertCustomerDetails(customerDetailsIns)
	if err != nil {
		resp = "Error while saving customer detail"
		renderJson(rw, http.StatusUnauthorized, resp)
		return
	}

	res := struct {
		ID      bson.ObjectId `json:"id"`
		Message string        `json:"message"`
	}{
		customerDetails.Id,
		`customer details saved successfully `,
	}

	renderJson(rw, http.StatusOK, res)
	return

}

func (p *Provider) GetCustomers(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	customerDetails, err := datastore.NewCustomerDetails(dbSession).FindAllCustomer()
	if err != nil {
		log.Printf("ERROR: GetCustomerDetailList %s", err)
		return
	}

	resp := struct {
		CustomerDetails *[]models.CustomerDetails `json:"customerDetails"`
	}{
		customerDetails,
	}

	renderJson(rw, http.StatusOK, resp)

}

func (p *Provider) UpdateCustomerDetails(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	customerID := chi.URLParam(r, "id")

	if !isObjectIDValid(customerID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid customer ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid customer ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	customerDetails, err := datastore.NewCustomerDetails(dbSession).FindByID(bson.ObjectIdHex(customerID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find customer details!.",
		}
		renderError(rw, err)
		return
	}

	if !parseJson(rw, r.Body, &customerDetails) {
		return
	}

	hasErr, validationErr := internal.CustomerValidation(p.db, *customerDetails)
	if hasErr {
		log.Printf("ERROR: InsertCustomerDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	customerDetails.LastUpdated = time.Now().UTC()

	if err = datastore.NewCustomerDetails(dbSession).Update(customerDetails.Id, *customerDetails); err != nil {
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
		customerDetails.Id,
		"Customer Details updated successfully.",
	}

	renderJson(rw, http.StatusOK, res)

}

func (p *Provider) RemoveCustomerDetails(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	customerID := chi.URLParam(r, "id")
	if !isObjectIDValid(customerID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid customer ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid customer ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	customerDetails, err := datastore.NewCustomerDetails(dbSession).FindByID(bson.ObjectIdHex(customerID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find customer details!.",
		}
		renderError(rw, err)
		return
	}

	err = internal.DeleteCustomerDetails(p.db, customerID, "customerDetails", customerDetails)
	if err != nil {
		p.log.Printf("ERROR: Handler - Delete - %q\n", err)
		renderError(rw, err)
		return
	}

	// Constructing response for client
	res := struct {
		Message string `json:"message"`
	}{
		"Customer Details Removed Successfully.",
	}

	renderJson(rw, http.StatusOK, res)
}
