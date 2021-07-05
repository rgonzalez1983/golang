package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"sync"
)

// MongoConnection
type MongoConnection struct {
	*mgo.Session
	m sync.Mutex
}

var BdName = os.Getenv("MONGO_DATABASE")

//NewConnection
func NewConnection(info *mgo.DialInfo) (*MongoConnection, error) {

	sess, err := mgo.DialWithInfo(info)

	if err != nil || sess == nil {
		return nil, err
	}

	sess.SetMode(mgo.Monotonic, true)
	sess.SetSafe(&mgo.Safe{})

	return &MongoConnection{sess, sync.Mutex{}}, nil
}

// Deleting Data
func (con *MongoConnection) DeleteData(collection string, id string) error {

	con.m.Lock()
	defer con.m.Unlock()

	c := con.DB(BdName).C(collection)
	oid := bson.ObjectIdHex(id)
	err := c.RemoveId(oid)
	return err
}

//Counting Data
func (con *MongoConnection) CountData(collection string, find string) (col int, err error) {

	con.m.Lock()
	defer con.m.Unlock()
	c := con

	if find != "" {
		col, err = c.DB(BdName).C(collection).Find(bson.M{}).Count()
	} else {
		col, err = c.DB(BdName).C(collection).Count()
	}
	return col, err
}

//Finding Data
func (con *MongoConnection) GetFindData(collection string, query bson.M, selector bson.M, fieldSort string, orderSort string) ([]interface{}, error) {

	con.m.Lock()
	defer con.m.Unlock()
	result := make([]interface{}, 0)
	err := con.DB(BdName).C(collection).Find(query).Select(selector).Sort(orderSort + fieldSort).All(&result)
	return result, err
}

//Inserting Data
func (con *MongoConnection) InsertData(collection string, ui interface{}) (err error) {

	con.m.Lock()
	defer con.m.Unlock()
	c := con

	col := c.DB(BdName).C(collection)
	event := ui

	err = col.Insert(&event)
	return err
}

//Updating Data
func (con *MongoConnection) UpdateData(collection string, ui interface{}, updt interface{}) (err error) {

	con.m.Lock()
	defer con.m.Unlock()
	c := con

	col := c.DB(BdName).C(collection)
	event := ui
	update := updt

	err = col.Update(&event, &update)
	return err
}
