package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type RawMaterial struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name" bson:"name"`
	Vendors     Vendor        `json:"vendors" bson:"vendors"`
	Grade       string        `json:"grade" bson:"grade"`
	Price       int           `json:"price" bson:"price"`
	DateCreated time.Time     `bson:"dateCreated" json:"dateCreated"`
	LastUpdated time.Time     `bson:"lastUpdated" json:"lastUpdated"`
}

type Vendor struct {
	VendorCode bson.ObjectId `json:"vendorCode" bson:"vendorCode"`
	VendorName string        `json:"vendorName" bson:"vendorName"`
}
