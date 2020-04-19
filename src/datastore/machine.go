package datastore

import (
	"angle/src/models"
	"angle/src/services"
	"angle/src/utils"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type machineDetails struct {
	// session    *mgo.Session
	collection *mgo.Collection
	// config     *config.Config
}

// NewUser creates a new user datastore
func NewMachineDetails(db *mgo.Session) services.MachineDetails {
	return &machineDetails{
		// session:    db,
		collection: db.DB("angle").C("machineDetails"),
		// config:     cfg,
	}
}

func (machine *machineDetails) InsertMachineDetails(machineDetail models.MachineDetails) (*models.MachineDetails, error) {

	if err := machine.collection.Insert(&machineDetail); err != nil {
		log.Printf("ERROR: InsertMachineDetails(%s) - %q\n", machineDetail.Name, err)
		return nil, err
	}

	return &machineDetail, nil

}

func (machine *machineDetails) Validate(machineDetail models.MachineDetails) (bool, map[string]interface{}) {

	v := &utils.Validation{}

	res := v.Required(machineDetail.Name).Key("name").Message("Enter name")
	if res.Ok {
		v.MaxSize(machineDetail.Name, 50).Key("name").Message("Name should not be more than 50 characters")
	}

	v.Required(machineDetail.Code).Key("code").Message("Enter code")
	v.Required(machineDetail.Make).Key("make").Message("Enter make")

	return v.HasErrors(), v.ErrorMap()
}

func (machine *machineDetails) FindAllMachineDetails() (*[]models.MachineDetails, error) {

	machineDetail := &[]models.MachineDetails{}
	err := machine.collection.Find(bson.M{}).All(machineDetail)
	if err != nil {
		log.Printf("ERROR: FindMachineDetails - %q\n", err)
		return nil, err
	}

	return machineDetail, nil
}

func (machine *machineDetails) FindByID(id bson.ObjectId) (*models.MachineDetails, error) {

	machineDetailsIns := &models.MachineDetails{}
	err := machine.collection.FindId(id).One(machineDetailsIns)
	if err != nil && err != mgo.ErrNotFound {
		log.Printf("ERROR: FindByID(%s) - %s\n", id, err)
		return nil, err
	}
	return machineDetailsIns, nil

}

func (machine *machineDetails) Update(id bson.ObjectId, machineDet models.MachineDetails) error {
	err := machine.collection.UpdateId(id, &machineDet)
	if err != nil {
		log.Printf("ERROR: Update(%s, %s) - %s\n", id, machineDet.Id, err)
		return err
	}
	return nil
}
