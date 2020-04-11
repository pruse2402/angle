package models

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type EmployeeDetails struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name" bson:"name"`
	Code        string        `json:"code" bson:"code"`
	Email       string        `json:"email" bson:"email"`
	PhoneNumber string        `json:"phoneNumber" bson:"phoneNumber"`
	DateCreated time.Time     `bson:"dateCreated" json:"dateCreated"`
	LastUpdated time.Time     `bson:"lastUpdated" json:"lastUpdated"`
}

type EmployeeDetailsByName []EmployeeDetails

func (a EmployeeDetailsByName) Len() int {
	return len(a)
}

func (a EmployeeDetailsByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a EmployeeDetailsByName) Less(i, j int) bool {
	return strings.ToUpper(a[i].Name) < strings.ToUpper(a[j].Name)
}
