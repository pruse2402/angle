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

func (p *Provider) InsertMachine(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	MachineDetailsIns := models.MachineDetails{}

	if ok := parseJson(rw, r.Body, &MachineDetailsIns); !ok {
		return
	}

	hasErr, validationErr := internal.MachineValidation(p.db, MachineDetailsIns)
	if hasErr {
		log.Printf("ERROR: InsertMachineDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	MachineDetailsIns.Id = bson.NewObjectId()
	MachineDetailsIns.DateCreated = time.Now().UTC()
	MachineDetailsIns.LastUpdated = time.Now().UTC()

	var resp string
	machineDetails, err := datastore.NewMachineDetails(dbSession).InsertMachineDetails(MachineDetailsIns)
	if err != nil {
		resp = "Error while saving machine detail"
		renderJson(rw, http.StatusUnauthorized, resp)
		return
	}

	res := struct {
		ID      bson.ObjectId `json:"id"`
		Message string        `json:"message"`
	}{
		machineDetails.Id,
		`machine details saved successfully `,
	}

	renderJson(rw, http.StatusOK, res)
	return

}

func (p *Provider) GetMachine(rw http.ResponseWriter, r *http.Request) {
	dbSession := p.db.Copy()
	defer dbSession.Close()

	machineDetails, err := datastore.NewMachineDetails(dbSession).FindAllMachineDetails()
	if err != nil {
		log.Printf("ERROR: GetMachineDetailList %s", err)
		return
	}

	sort.Sort(models.MachineDetailsByName(*machineDetails))

	resp := struct {
		MachineDetails *[]models.MachineDetails `json:"machineDetails"`
	}{
		machineDetails,
	}

	renderJson(rw, http.StatusOK, resp)

}

func (p *Provider) UpdateMachine(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	machineID := chi.URLParam(r, "id")

	if !isObjectIDValid(machineID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid machine ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid machine ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	machineDetail, err := datastore.NewMachineDetails(dbSession).FindByID(bson.ObjectIdHex(machineID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find machine details!.",
		}
		renderError(rw, err)
		return
	}

	if !parseJson(rw, r.Body, &machineDetail) {
		return
	}

	hasErr, validationErr := internal.MachineValidation(p.db, *machineDetail)
	if hasErr {
		log.Printf("ERROR: UpdateMachineDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	machineDetail.LastUpdated = time.Now().UTC()

	if err = datastore.NewMachineDetails(dbSession).Update(machineDetail.Id, *machineDetail); err != nil {
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
		machineDetail.Id,
		"machine detail updated successfully.",
	}

	renderJson(rw, http.StatusOK, res)
}

func (p *Provider) RemoveMachine(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	machineID := chi.URLParam(r, "id")
	if !isObjectIDValid(machineID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid Machine ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid Machine ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	machineDetails, err := datastore.NewMachineDetails(dbSession).FindByID(bson.ObjectIdHex(machineID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find machine details!.",
		}
		renderError(rw, err)
		return
	}

	err = internal.DeleteMachine(p.db, machineID, "machineDetails", machineDetails)
	if err != nil {
		p.log.Printf("ERROR: Handler - Delete - %q\n", err)
		renderError(rw, err)
		return
	}

	// Constructing response for client
	res := struct {
		Message string `json:"message"`
	}{
		"Machine Details Removed Successfully.",
	}

	renderJson(rw, http.StatusOK, res)
}
