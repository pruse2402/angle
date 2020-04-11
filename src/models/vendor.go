package models

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type VendorDetails struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name" bson:"name"`
	Address     string        `json:"address" bson:"address"`
	Pincode     string        `json:"pincode" bson:"pincode"`
	GstIn       string        `json:"gstIn" bson:"gstIn"`
	DateCreated time.Time     `bson:"dateCreated" json:"dateCreated"`
	LastUpdated time.Time     `bson:"lastUpdated" json:"lastUpdated"`
}

type VendorDetailsByName []VendorDetails

func (a VendorDetailsByName) Len() int {
	return len(a)
}

func (a VendorDetailsByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a VendorDetailsByName) Less(i, j int) bool {
	return strings.ToUpper(a[i].Name) < strings.ToUpper(a[j].Name)
}
