package entity

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Person struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string        `json:"name"`
	LastName string        `json:"lastname"`
	CI       string        `json:"ci"`
	Country  string        `json:"country,omitempty"`
	Address  string        `json:"address,omitempty"`
	Age      int           `json:"age,omitempty"`
	Created  time.Time     `json:"created"`
	Updated  time.Time     `json:"updated"`
}
