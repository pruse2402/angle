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

	resp := struct {
		EmployeeDetails *[]models.EmployeeDetails `json:"employeeDetails"`
	}{
		employeeDetails,
	}

	renderJson(rw, http.StatusOK, resp)

}
