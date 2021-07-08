package persistance

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"go_project/db"
	"go_project/internal"
	"go_project/internal/entity"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"time"
)

type PersonRepository interface {
	CreatePerson(r *interface{}) (template interface{}, error error, status int)
	UpdatePerson(id string, r *interface{}) (template interface{}, error error, status int)
	ListPersons() (templates []interface{}, error error, status int)
	GetPerson(id string) (template interface{}, error error, status int)
	DeletePerson(id string) (template interface{}, error error, status int)
}

type personRepository struct {
	connMgo *db.MongoConnection
}

func NewPersonRepository(connMgo *db.MongoConnection) PersonRepository {
	_ = connMgo.EnsureIndex(internal.CollectionPerson, []string{internal.FieldCi})
	return &personRepository{
		connMgo: connMgo,
	}
}

//Creating Person
func (p personRepository) CreatePerson(r *interface{}) (template interface{}, error error, status int) {
	object := p.ToEntityObject(*r)
	object.Created, object.Updated = time.Now(), time.Now()
	err := p.connMgo.InsertData(internal.CollectionPerson, object)
	if err != nil {
		return object, errors.New(internal.MsgResponseObjectExists), http.StatusConflict
	}
	return object, nil, http.StatusCreated
}

//Updating Person
func (p personRepository) UpdatePerson(id string, r *interface{}) (template interface{}, error error, status int) {
	objectNew := p.ToEntityObject(*r)
	objectNew.Updated = time.Now()
	filter := bson.M{internal.Field__id: bson.ObjectIdHex(id)}
	found, _ := p.GetFindPersons(internal.CollectionPerson, filter, bson.M{}, internal.FieldUpdated, internal.OrderDesc)
	if len(found) > 0 {
		document, _ := p.ToDocument(objectNew)
		update := bson.M{internal.MongoDB__set: *document}
		err := p.connMgo.UpdateData(internal.CollectionPerson, filter, update)
		if err != nil {
			return objectNew, errors.New(internal.MsgResponseObjectExists), http.StatusConflict
		}
		return objectNew, nil, http.StatusCreated
	}
	return objectNew, errors.New(internal.MsgResponseServerError), http.StatusInternalServerError
}

//Listing Persons
func (p personRepository) ListPersons() (templates []interface{}, error error, status int) {
	collection := internal.CollectionPerson
	found, _ := p.GetFindPersons(collection, bson.M{}, bson.M{}, internal.FieldLastname, internal.OrderAsc)
	return found, nil, http.StatusCreated
}

//Getting Person
func (p personRepository) GetPerson(id string) (template interface{}, error error, status int) {
	collection := internal.CollectionPerson
	filter := bson.M{internal.Field__id: bson.ObjectIdHex(id)}
	//filter, _ := p.ToDocument(*r)
	found, _ := p.GetFindPersons(collection, filter, bson.M{}, internal.FieldLastname, internal.OrderAsc)
	if len(found) == 0 {
		return internal.ValueEmpty, errors.New(internal.MsgResponseNoData), http.StatusInternalServerError
	}
	return found[0], nil, http.StatusCreated
}

//Deleting Person
func (p personRepository) DeletePerson(id string) (template interface{}, error error, status int) {
	collection := internal.CollectionPerson
	filter := bson.M{internal.Field__id: bson.ObjectIdHex(id)}
	found, _ := p.GetFindPersons(collection, filter, bson.M{}, internal.FieldUpdated, internal.OrderDesc)
	if len(found) > 0 {
		_ = p.connMgo.DeleteData(collection, id)
		return found[0], nil, http.StatusCreated
	}
	return nil, errors.New(internal.MsgResponseNoData), http.StatusInternalServerError
}

//FUNCIONES AUXILIARES

func (p personRepository) GetFindPersons(collection string, query bson.M, selector bson.M, fieldSort string, orderSort string) (items []interface{}, err error) {
	items, err = p.connMgo.GetFindData(collection, query, selector, fieldSort, orderSort)
	return items, err
}

func (p personRepository) ToDocument(v interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func (p personRepository) ToEntityObject(i interface{}) entity.Person {
	person := entity.Person{}
	if reflect.TypeOf(i).String() != "bson.M" {
		m := i.(map[string]interface{})
		_ = mapstructure.Decode(m, &person)
	} else {
		bsonBytes, _ := bson.Marshal(i)
		_ = bson.Unmarshal(bsonBytes, &person)
	}
	return person
}
