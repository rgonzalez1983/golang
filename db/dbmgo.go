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

//Create Connection
func (con *MongoConnection) CreateConnection() *mgo.Database {
	con.m.Lock()
	defer con.m.Unlock()
	return con.DB(BdName)
}

// Deleting Data
func (con *MongoConnection) DeleteData(collection string, id string) error {
	c := con.CreateConnection().C(collection)
	oid := bson.ObjectIdHex(id)
	err := c.RemoveId(oid)
	return err
}

//Counting Data
func (con *MongoConnection) CountData(collection string, find string) (col int, err error) {
	c := con.CreateConnection().C(collection)
	if find != "" {
		col, err = c.Find(bson.M{}).Count()
	} else {
		col, err = c.Count()
	}
	return col, err
}

//Finding Data
func (con *MongoConnection) GetFindData(collection string, query bson.M, selector bson.M, fieldSort string, orderSort string) ([]interface{}, error) {
	c := con.CreateConnection().C(collection)
	result := make([]interface{}, 0)
	err := c.Find(query).Select(selector).Sort(orderSort + fieldSort).All(&result)
	return result, err
}

//Inserting Data
func (con *MongoConnection) InsertData(collection string, ui interface{}) (err error) {
	c := con.CreateConnection().C(collection)
	event := ui
	err = c.Insert(&event)
	return err
}

//Updating Data
func (con *MongoConnection) UpdateData(collection string, ui interface{}, updt interface{}) (err error) {
	c := con.CreateConnection().C(collection)
	event := ui
	update := updt
	err = c.Update(&event, &update)
	return err
}

//Indexes
func (con *MongoConnection) EnsureIndex(collection string, arrayIndexes []string) (err error) {
	c := con.CreateConnection().C(collection)
	for _, key := range arrayIndexes {
		index := mgo.Index{
			Key:    []string{key},
			Unique: true,
		}
		if err := c.EnsureIndex(index); err != nil {
			return err
		}
	}
	return err
}
