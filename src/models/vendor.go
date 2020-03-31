package models

import (
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
