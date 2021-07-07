package entity

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Person struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name,omitempty" bson:"name,omitempty"`
	LastName string        `json:"lastname,omitempty" bson:"lastname,omitempty"`
	CI       string        `json:"ci,omitempty" bson:"ci,omitempty"`
	Country  string        `json:"country,omitempty" bson:"country,omitempty"`
	Address  string        `json:"address,omitempty" bson:"address,omitempty"`
	Age      int           `json:"age,omitempty" bson:"age,omitempty"`
	Created  time.Time     `json:"created"`
	Updated  time.Time     `json:"updated"`
}
