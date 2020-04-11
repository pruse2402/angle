package datastore

import (
	"angle/src/models"
	"angle/src/services"
	"angle/src/utils"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (employee *employeeDetails) FindAllEmployee() (*[]models.EmployeeDetails, error) {

	employeeDetail := &[]models.EmployeeDetails{}
	err := employee.collection.Find(bson.M{}).All(employeeDetail)
	if err != nil {
		log.Printf("ERROR: FindEmployeeDetails - %q\n", err)
		return nil, err
	}

	return employeeDetail, nil
}

func (employee *employeeDetails) FindByID(id bson.ObjectId) (*models.EmployeeDetails, error) {

	employeeDetailsIns := &models.EmployeeDetails{}
	err := employee.collection.FindId(id).One(employeeDetailsIns)
	if err != nil && err != mgo.ErrNotFound {
		log.Printf("ERROR: FindByID(%s) - %s\n", id, err)
		return nil, err
	}
	return employeeDetailsIns, nil

}

func (employee *employeeDetails) Update(id bson.ObjectId, empDetail models.EmployeeDetails) error {
	err := employee.collection.UpdateId(id, &empDetail)
	if err != nil {
		log.Printf("ERROR: Update(%s, %s) - %s\n", id, empDetail.Id, err)
		return err
	}
	return nil
}
