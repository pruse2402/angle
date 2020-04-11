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

func (p *Provider) InsertEmployee(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	EmployeeDetailsIns := models.EmployeeDetails{}

	if ok := parseJson(rw, r.Body, &EmployeeDetailsIns); !ok {
		return
	}

	hasErr, validationErr := internal.EmployeeValidation(p.db, EmployeeDetailsIns)
	if hasErr {
		log.Printf("ERROR: InsertEmployeeDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	EmployeeDetailsIns.Id = bson.NewObjectId()
	EmployeeDetailsIns.DateCreated = time.Now().UTC()
	EmployeeDetailsIns.LastUpdated = time.Now().UTC()

	var resp string
	employeeDetails, err := datastore.NewEmployeeDetails(dbSession).InsertEmployeeDetails(EmployeeDetailsIns)
	if err != nil {
		resp = "Error while saving employee detail"
		renderJson(rw, http.StatusUnauthorized, resp)
		return
	}

	res := struct {
		ID      bson.ObjectId `json:"id"`
		Message string        `json:"message"`
	}{
		employeeDetails.Id,
		`employee details saved successfully `,
	}

	renderJson(rw, http.StatusOK, res)
	return

}

func (p *Provider) GetEmployee(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	employeeDetails, err := datastore.NewEmployeeDetails(dbSession).FindAllEmployee()
	if err != nil {
		log.Printf("ERROR: GetEmployeeDetailList %s", err)
		return
	}

	sort.Sort(models.EmployeeDetailsByName(*employeeDetails))

	resp := struct {
		EmployeeDetails *[]models.EmployeeDetails `json:"employeeDetails"`
	}{
		employeeDetails,
	}

	renderJson(rw, http.StatusOK, resp)

}

func (p *Provider) UpdateEmployee(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	employeeID := chi.URLParam(r, "id")

	if !isObjectIDValid(employeeID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid employee ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid employee ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	employeeDetail, err := datastore.NewEmployeeDetails(dbSession).FindByID(bson.ObjectIdHex(employeeID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find employee details!.",
		}
		renderError(rw, err)
		return
	}

	if !parseJson(rw, r.Body, &employeeDetail) {
		return
	}

	hasErr, validationErr := internal.EmployeeValidation(p.db, *employeeDetail)
	if hasErr {
		log.Printf("ERROR: UpdateEmployeeDetails - %q\n", validationErr)
		err := &errs.AppError{
			Message: "Validation Error(s)",
			Errors:  validationErr,
		}
		respondError(rw, http.StatusBadRequest, err)
		return
	}

	employeeDetail.LastUpdated = time.Now().UTC()

	if err = datastore.NewEmployeeDetails(dbSession).Update(employeeDetail.Id, *employeeDetail); err != nil {
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
		employeeDetail.Id,
		"raw materials updated successfully.",
	}

	renderJson(rw, http.StatusOK, res)
}

func (p *Provider) RemoveEmployee(rw http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	employeeID := chi.URLParam(r, "id")
	if !isObjectIDValid(employeeID) {
		p.log.Printf("ERROR: Handler - Update - %q\n", `Invalid Employee ID Supplied.`)
		err := &errs.UIErr{
			Code:    http.StatusBadRequest,
			Message: "Invalid Employee ID Supplied.",
		}
		renderError(rw, err)
		return
	}

	employeeDetails, err := datastore.NewEmployeeDetails(dbSession).FindByID(bson.ObjectIdHex(employeeID))
	if err != nil {
		p.log.Printf("ERROR: Handler - Update - %q\n", err)
		err := &errs.UIErr{
			Code:    http.StatusNotFound,
			Message: "Could not find employee details!.",
		}
		renderError(rw, err)
		return
	}

	err = internal.DeleteEmployee(p.db, employeeID, "employeeDetails", employeeDetails)
	if err != nil {
		p.log.Printf("ERROR: Handler - Delete - %q\n", err)
		renderError(rw, err)
		return
	}

	// Constructing response for client
	res := struct {
		Message string `json:"message"`
	}{
		"Employee Details Removed Successfully.",
	}

	renderJson(rw, http.StatusOK, res)
}
