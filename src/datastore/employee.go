package datastore

import (
	"angle/src/models"
	"angle/src/services"
	"angle/src/utils"
	"log"

	"gopkg.in/mgo.v2"
)

type employeeDetails struct {
	// session    *mgo.Session
	collection *mgo.Collection
	// config     *config.Config
}

// NewUser creates a new user datastore
func NewEmployeeDetails(db *mgo.Session) services.EmployeeDetails {
	return &employeeDetails{
		// session:    db,
		collection: db.DB("angle").C("employeeDetails"),
		// config:     cfg,
	}
}

func (employee *employeeDetails) InsertEmployeeDetails(employeeDetail models.EmployeeDetails) (*models.EmployeeDetails, error) {

	if err := employee.collection.Insert(&employeeDetail); err != nil {
		log.Printf("ERROR: InsertEmployeeDetails(%s) - %q\n", employeeDetail.Name, err)
		return nil, err
	}

	return &employeeDetail, nil

}

func (employee *employeeDetails) Validate(employeeDetail models.EmployeeDetails) (bool, map[string]interface{}) {

	v := &utils.Validation{}

	res := v.Required(employeeDetail.Name).Key("name").Message("Enter name")
	if res.Ok {
		v.MaxSize(employeeDetail.Name, 50).Key("name").Message("Name should not be more than 50 characters")
	}

	if employeeDetail.PhoneNumber != "" {
		v.MinSize(employeeDetail.PhoneNumber, 10).Key("phoneNumber").Message("Enter valid phone number")
	}

	res = v.Required(employeeDetail.Code).Key("code").Message("Enter code")
	if res.Ok {
		v.MaxSize(employeeDetail.Code, 10).Key("code").Message("Code should not be more than 10 characters")

	}

	return v.HasErrors(), v.ErrorMap()
}
