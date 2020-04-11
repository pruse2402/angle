package models

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type CustomerDetails struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name" bson:"name"`
	Address     string        `json:"address" bson:"address"`
	Pincode     string        `json:"pincode" bson:"pincode"`
	GstIn       string        `json:"gstIn" bson:"gstIn"`
	DateCreated time.Time     `bson:"dateCreated" json:"dateCreated"`
	LastUpdated time.Time     `bson:"lastUpdated" json:"lastUpdated"`
}

type CustomerDetailsByName []CustomerDetails

func (a CustomerDetailsByName) Len() int {
	return len(a)
}

func (a CustomerDetailsByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a CustomerDetailsByName) Less(i, j int) bool {
	return strings.ToUpper(a[i].Name) < strings.ToUpper(a[j].Name)
}
