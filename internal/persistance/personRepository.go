package persistance

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"go_project/db"
	"go_project/internal/entity"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"time"
)

type PersonRepository interface {
	CreatePerson(r *interface{}) (template interface{}, error error, status int)
	UpdatePerson(r *interface{}) (template interface{}, error error, status int)
	ListPersons() (templates []interface{}, error error, status int)
	GetPerson(r *interface{}) (template interface{}, error error, status int)
	DeletePerson(r *interface{}) (template interface{}, error error, status int)
}

type personRepository struct {
	connMgo *db.MongoConnection
}

func NewPersonRepository(connMgo *db.MongoConnection) PersonRepository {
	return &personRepository{
		connMgo: connMgo,
	}
}

func (p personRepository) CreatePerson(r *interface{}) (template interface{}, error error, status int) {
	object := p.ToEntityObject(*r)
	collection, filter := "person", bson.M{"ci": object.CI}
	object.Created, object.Updated = time.Now(), time.Now()
	found, _ := p.GetFindPersons(collection, filter, bson.M{}, "updated", "-")
	if len(found) > 0 {
		return object, errors.New("PERSONA EXISTENTE"), http.StatusNotAcceptable
	}
	_ = p.connMgo.InsertData(collection, object)
	return object, nil, http.StatusAccepted
}

func (p personRepository) UpdatePerson(r *interface{}) (template interface{}, error error, status int) {
	objectNew := p.ToEntityUpdateObject(*r)
	collection, filter := "person", bson.M{"ci": objectNew.Values.CI}
	objectNew.Values.Updated = time.Now()
	found, _ := p.GetFindPersons(collection, filter, bson.M{}, "updated", "-")
	if len(found) > 0 {
		objectOld := p.ToEntityObject(found[0])
		if objectOld.ID.String() != objectNew.ID {
			return objectNew.Values, errors.New("PERSONA EXISTENTE"), http.StatusNotAcceptable
		}
	}
	filter = bson.M{"_id": bson.ObjectIdHex(objectNew.ID)}
	found, _ = p.GetFindPersons(collection, filter, bson.M{}, "updated", "-")
	if len(found) > 0 {
		document, _ := p.ToDocument(objectNew.Values)
		update := bson.M{"$set": *document}
		_ = p.connMgo.UpdateData(collection, filter, update)
		return objectNew, nil, http.StatusAccepted
	}
	return objectNew, errors.New("ERROR DE SERVIDOR"), http.StatusServiceUnavailable
}

func (p personRepository) ListPersons() (templates []interface{}, error error, status int) {
	collection := "person"
	found, _ := p.GetFindPersons(collection, bson.M{}, bson.M{}, "lastname", "")
	return found, nil, http.StatusAccepted
}

func (p personRepository) GetPerson(r *interface{}) (template interface{}, error error, status int) {
	collection := "person"
	filter, _ := p.ToDocument(*r)
	found, _ := p.GetFindPersons(collection, *filter, bson.M{}, "lastname", "")
	if len(found) == 0 {
		return "", errors.New("SIN DATOS EXISTENTES"), http.StatusNotFound
	}
	return found[0], nil, http.StatusAccepted
}

func (p personRepository) DeletePerson(r *interface{}) (template interface{}, error error, status int) {
	objectDelete := p.ToEntityDeleteObject(*r)
	collection, filter := "person", bson.M{"_id": bson.ObjectIdHex(objectDelete.ID)}
	found, _ := p.GetFindPersons(collection, filter, bson.M{}, "updated", "-")
	if len(found) > 0 {
		_ = p.connMgo.DeleteData(collection, objectDelete.ID)
		return found[0], nil, http.StatusAccepted
	}
	return nil, errors.New("SIN DATOS EXISTENTES"), http.StatusServiceUnavailable
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

func (p personRepository) ToEntityUpdateObject(i interface{}) entity.UpdatePerson {
	m, personUpdate := i.(map[string]interface{}), entity.UpdatePerson{}
	_ = mapstructure.Decode(m, &personUpdate)
	return personUpdate
}

func (p personRepository) ToEntityDeleteObject(i interface{}) entity.DeletePerson {
	m, personDelete := i.(map[string]interface{}), entity.DeletePerson{}
	_ = mapstructure.Decode(m, &personDelete)
	return personDelete
}
