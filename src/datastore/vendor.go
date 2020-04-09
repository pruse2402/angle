package datastore

import (
	"angle/src/models"
	"angle/src/services"
	"angle/src/utils"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// user is holds information for database
type vendorDetails struct {
	// session    *mgo.Session
	collection *mgo.Collection
	// config     *config.Config
}

// NewUser creates a new user datastore
func NewVendorDetails(db *mgo.Session) services.VendorDetails {
	return &vendorDetails{
		// session:    db,
		collection: db.DB("angle").C("vendorDetails"),
		// config:     cfg,
	}
}

func (vendor *vendorDetails) InsertVendorDetails(vendorDetail models.VendorDetails) (*models.VendorDetails, error) {

	if err := vendor.collection.Insert(&vendorDetail); err != nil {
		log.Printf("ERROR: InsertVendorDetails(%s) - %q\n", vendorDetail.Name, err)
		return nil, err
	}

	return &vendorDetail, nil

}

func (vendor *vendorDetails) Validate(vendorDetail models.VendorDetails) (bool, map[string]interface{}) {

	v := &utils.Validation{}

	//Escape special characters to avoid duplicates
	vendorDetailsGstIn := utils.AddEscapeString(vendorDetail.GstIn)

	res := v.Required(vendorDetail.Address).Key("address").Message("Enter address")

	res = v.Required(vendorDetail.Name).Key("name").Message("Enter name")
	if res.Ok {
		res = v.MaxSize(vendorDetail.Name, 50).Key("name").Message("Name should not be more than 50 characters")
		if res.Ok {
			query := bson.M{"gstIn": bson.RegEx{`^` + vendorDetailsGstIn + `$`, `i`}}
			if vendorDetail.Id.Hex() != "" {
				query["_id"] = bson.M{"$ne": vendorDetail.Id}
			}

			count, _ := vendor.collection.Find(query).Count()
			if count > 0 {
				v.Error("GSTIN already exists").Key("gstIn")
			}
		}
	}

	res = v.Required(vendorDetail.Address).Key("address").Message("Enter address")
	if res.Ok {
		v.Check(vendorDetail.Pincode, utils.Required{}, utils.MinSize{4}, utils.MaxSize{6}).Key("pincode").Message("Enter valid pin code")
	}

	return v.HasErrors(), v.ErrorMap()
}

func (vendor *vendorDetails) FindAllVendor() (*[]models.VendorDetails, error) {

	vendorDetails := &[]models.VendorDetails{}
	err := vendor.collection.Find(bson.M{}).All(vendorDetails)
	if err != nil {
		log.Printf("ERROR: FindVendorDetails - %q\n", err)
		return nil, err
	}

	return vendorDetails, nil
}

func (vendor *vendorDetails) FindByID(id bson.ObjectId) (*models.VendorDetails, error) {

	vendorDetailsIns := &models.VendorDetails{}
	err := vendor.collection.FindId(id).One(vendorDetailsIns)
	if err != nil && err != mgo.ErrNotFound {
		log.Printf("ERROR: FindByID(%s) - %s\n", id, err)
		return nil, err
	}
	return vendorDetailsIns, nil

}

func (vendor *vendorDetails) Update(id bson.ObjectId, vendorDetail models.VendorDetails) error {
	err := vendor.collection.UpdateId(id, &vendorDetail)
	if err != nil {
		log.Printf("ERROR: Update(%s, %s) - %s\n", id, vendorDetail.Id, err)
		return err
	}
	return nil
}
