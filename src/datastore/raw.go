package datastore

import (
	"angle/src/models"
	"angle/src/services"
	"angle/src/utils"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type rawMaterial struct {
	// session    *mgo.Session
	collection *mgo.Collection
	// config     *config.Config
}

// NewUser creates a new user datastore
func NewRawMaterial(db *mgo.Session) services.RawMaterial {
	return &rawMaterial{
		// session:    db,
		collection: db.DB("angle").C("rawMaterial"),
		// config:     cfg,
	}
}

func (raw *rawMaterial) InsertRawMaterial(rawMaterial models.RawMaterial) (*models.RawMaterial, error) {

	if err := raw.collection.Insert(&rawMaterial); err != nil {
		log.Printf("ERROR: InsertMenuItem(%s) - %q\n", rawMaterial.Name, err)
		return nil, err
	}

	return &rawMaterial, nil

}

func (raw *rawMaterial) Validate(rawMaterial models.RawMaterial) (bool, map[string]interface{}) {

	v := &utils.Validation{}

	//Escape special characters to avoid duplicates
	//rawMaterialName := utils.AddEscapeString(rawMaterial.Name)

	res := v.Required(rawMaterial.Grade).Key("grade").Message("Enter grade")

	res = v.Required(rawMaterial.Name).Key("name").Message("Enter name")
	if res.Ok {
		v.MaxSize(rawMaterial.Name, 50).Key("name").Message("Name should not be more than 50 characters")
	}

	// if res.Ok {
	// 	res = v.MaxSize(rawMaterial.Name, 50).Key("name").Message("Name should not be more than 50 characters")
	// 	if res.Ok {
	// 		query := bson.M{"gstIn": bson.RegEx{`^` + rawMaterialName + `$`, `i`}}
	// 		if rawMaterial.Id.Hex() != "" {
	// 			query["_id"] = bson.M{"$ne": rawMaterial.Id}
	// 		}

	// 		count, _ := raw.collection.Find(query).Count()
	// 		if count > 0 {
	// 			v.Error("Name already exists").Key("name")
	// 		}
	// 	}
	// }

	res = v.Required(rawMaterial.Vendors.VendorName).Key("vendorName").Message("Enter vendor name")

	res = v.Required(rawMaterial.Vendors.VendorCode).Key("vendorCode").Message("Enter vendor code")

	// if res.Ok {
	// 	v.Check(customerDetail.Pincode, utils.Required{}, utils.MinSize{4}, utils.MaxSize{6}).Key("pincode").Message("Enter valid pin code")
	// }

	return v.HasErrors(), v.ErrorMap()
}

func (raw *rawMaterial) FindAllRawMaterial() (*[]models.RawMaterial, error) {

	rawMaterial := &[]models.RawMaterial{}
	err := raw.collection.Find(bson.M{}).All(rawMaterial)
	if err != nil {
		log.Printf("ERROR: FindRawMaterialDetails - %q\n", err)
		return nil, err
	}

	return rawMaterial, nil
}
