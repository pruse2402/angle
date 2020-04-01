package datastore

import (
	"angle/src/models"
	"angle/src/services"
	"angle/src/utils"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type customerDetails struct {
	// session    *mgo.Session
	collection *mgo.Collection
	// config     *config.Config
}

// NewUser creates a new user datastore
func NewCustomerDetails(db *mgo.Session) services.CustomerDetails {
	return &customerDetails{
		// session:    db,
		collection: db.DB("angle").C("customerDetails"),
		// config:     cfg,
	}
}

func (customer *customerDetails) InsertCustomerDetails(customerDetail models.CustomerDetails) (*models.CustomerDetails, error) {

	if err := customer.collection.Insert(&customerDetail); err != nil {
		log.Printf("ERROR: InsertMenuItem(%s) - %q\n", customerDetail.Name, err)
		return nil, err
	}

	return &customerDetail, nil

}

func (customer *customerDetails) Validate(customerDetail models.CustomerDetails) (bool, map[string]interface{}) {

	v := &utils.Validation{}

	//Escape special characters to avoid duplicates
	customerDetailsGstIn := utils.AddEscapeString(customerDetail.GstIn)

	res := v.Required(customerDetail.Address).Key("address").Message("Enter address")

	res = v.Required(customerDetail.Name).Key("name").Message("Enter name")
	if res.Ok {
		res = v.MaxSize(customerDetail.Name, 50).Key("name").Message("Name should not be more than 50 characters")
		if res.Ok {
			query := bson.M{"gstIn": bson.RegEx{`^` + customerDetailsGstIn + `$`, `i`}}
			if customerDetail.Id.Hex() != "" {
				query["_id"] = bson.M{"$ne": customerDetail.Id}
			}

			count, _ := customer.collection.Find(query).Count()
			if count > 0 {
				v.Error("GSTIN already exists").Key("gstIn")
			}
		}
	}

	res = v.Required(customerDetail.Address).Key("address").Message("Enter address")
	if res.Ok {
		v.Check(customerDetail.Pincode, utils.Required{}, utils.MinSize{4}, utils.MaxSize{6}).Key("pincode").Message("Enter valid pin code")
	}

	return v.HasErrors(), v.ErrorMap()
}

func (customer *customerDetails) FindAllCustomer() (*[]models.CustomerDetails, error) {

	customerDetails := &[]models.CustomerDetails{}
	err := customer.collection.Find(bson.M{}).All(customerDetails)
	if err != nil {
		log.Printf("ERROR: FindCustomerDetails - %q\n", err)
		return nil, err
	}

	return customerDetails, nil
}

func (customer *customerDetails) FindByID(id bson.ObjectId) (*models.CustomerDetails, error) {

	customerDetailsIns := &models.CustomerDetails{}
	err := customer.collection.FindId(id).One(customerDetailsIns)
	if err != nil && err != mgo.ErrNotFound {
		log.Printf("ERROR: FindByID(%s) - %s\n", id, err)
		return nil, err
	}
	return customerDetailsIns, nil

}

func (customer *customerDetails) Update(id bson.ObjectId, customerDetail models.CustomerDetails) error {
	err := customer.collection.UpdateId(id, &customerDetail)
	if err != nil {
		log.Printf("ERROR: Update(%s, %s) - %s\n", id, customerDetail.Id, err)
		return err
	}
	return nil
}
