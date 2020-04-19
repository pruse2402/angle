package models

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type MachineDetails struct {
	Id             bson.ObjectId `bson:"_id" json:"id"`
	Name           string        `json:"name" bson:"name"`
	Code           string        `json:"code" bson:"code"`
	Make           string        `json:"make" bson:"make"`
	MachineType    string        `json:"machineType" bson:"machineType"`
	Tonnage        string        `json:"tonnage" bson:"tonnage"`
	TieBarDistance int           `json:"tieBarDistance" bson:"tieBarDistance"`
	PlattenSizeLen int           `json:"plattenSizeLen" bson:"plattenSizeLen"`
	PlattenSizewid int           `json:"PlattenSizewid" bson:"PlattenSizewid"`
	MachineDimLen  int           `json:"machineDimLen" bson:"machineDimLen"`
	MachineDimWid  int           `json:"machineDimWid" bson:"machineDimWid"`
	MachineDimHt   int           `json:"machineDimHt" bson:"machineDimHt"`
	ScrewSize      int           `json:"screwSize" bson:"screwSize"`
	DateCreated    time.Time     `bson:"dateCreated" json:"dateCreated"`
	LastUpdated    time.Time     `bson:"lastUpdated" json:"lastUpdated"`
}

type MachineDetailsByName []MachineDetails

func (a MachineDetailsByName) Len() int {
	return len(a)
}

func (a MachineDetailsByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a MachineDetailsByName) Less(i, j int) bool {
	return strings.ToUpper(a[i].Name) < strings.ToUpper(a[j].Name)
}
