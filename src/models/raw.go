package models

import (
	"strings"
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

type RawMaterialByName []RawMaterial

func (a RawMaterialByName) Len() int {
	return len(a)
}

func (a RawMaterialByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a RawMaterialByName) Less(i, j int) bool {
	return strings.ToUpper(a[i].Name) < strings.ToUpper(a[j].Name)
}
