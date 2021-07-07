package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go_project/internal"
	"go_project/internal/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a middleware.App

func TestMain(m *testing.M) {
	a = App()
	code := m.Run()
	os.Exit(code)
}

func App() middleware.App {
	a := middleware.App{}
	_ = a.Initialize(internal.ValueEmpty, internal.ValueEmpty)
	return a
}

func ResponseToJSON(responseBody string) map[string]interface{} {
	var JSON map[string]interface{}
	_ = json.Unmarshal([]byte(responseBody), &JSON)
	return JSON
}

func TestListPersons(t *testing.T) {
	request, _ := http.NewRequest(internal.HTTP_POST, internal.URLListingAll, nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)
	values := []interface{}{internal.KeyType, internal.TEST, internal.KeyURL, internal.URLListingAll, internal.KeyMessage, internal.MsgResponseListingAll}
	a.LoggingOperation(values...)
	responseBody := ResponseToJSON(response.Body.String())
	length := len(responseBody[internal.KeyResponseData].([]interface{}))
	assert.Equal(t, 201, response.Code, "EXPECTED 201")
	assert.Equal(t, length, len(responseBody[internal.KeyResponseData].([]interface{})), "EXPECTED "+strconv.Itoa(length))
}

func TestGetPerson(t *testing.T) {
	payload := []byte(`{"age":37}`)
	request, _ := http.NewRequest(internal.HTTP_POST, internal.URLGettingOne, bytes.NewBuffer(payload))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)
	message := func() string {
		if response.Code == 500 {
			return internal.MsgResponseNoData
		} else {
			return internal.MsgResponseGettingOne
		}
	}()
	values := []interface{}{internal.KeyType, internal.TEST, internal.KeyURL, internal.URLGettingOne, internal.KeyMessage, message}
	a.LoggingOperation(values...)
	responseBody := ResponseToJSON(response.Body.String())
	assert.Equal(t, 201, response.Code, "EXPECTED 201")
	assert.Equal(t, message, responseBody[internal.KeyResponseMessage].(interface{}), "EXPECTED "+message)
	assert.Equal(t, "83110715463", responseBody[internal.KeyResponseData].(map[string]interface{})["ci"], "EXPECTED 83110715463")
}

func TestCreatePerson(t *testing.T) {
	payload := []byte(`{"name" : "JUAN BRAULIO",
						"lastname" : "HERNANDEZ NAPOLES",
						"ci" : "96092017065",
						"country" : "Cuba",
						"address" : "Buenos Aires, Camaguey",
						"age" : 24}`)
	request, _ := http.NewRequest(internal.HTTP_POST, internal.URLCreatingOne, bytes.NewBuffer(payload))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)
	message := func() string {
		if response.Code == 500 {
			return internal.MsgResponseObjectExists
		} else {
			return internal.MsgResponseCreatingOne
		}
	}()
	values := []interface{}{internal.KeyType, internal.TEST, internal.KeyURL, internal.URLCreatingOne, internal.KeyMessage, message}
	a.LoggingOperation(values...)
	responseBody := ResponseToJSON(response.Body.String())
	assert.Equal(t, 500, response.Code, "EXPECTED 500")
	assert.Equal(t, message, responseBody[internal.KeyResponseMessage].(interface{}), "EXPECTED "+message)
	assert.Contains(t, responseBody[internal.KeyResponseData].(map[string]interface{})["ci"], "96092017065", "EXPECTED 96092017065")
}

func TestUpdatePerson(t *testing.T) {
	payload := []byte(`{
						"id" : "60de3ce4befb1fb4a19d8dfd",
						"values" : {
							"name" : "ANA M."
                         }
						}`)
	request, _ := http.NewRequest(internal.HTTP_POST, internal.URLUpdatingOne, bytes.NewBuffer(payload))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)
	message := func() string {
		if response.Code == 500 {
			return internal.MsgResponseObjectExists
		} else {
			return internal.MsgResponseUpdatingOne
		}
	}()
	values := []interface{}{internal.KeyType, internal.TEST, internal.KeyURL, internal.URLUpdatingOne, internal.KeyMessage, message}
	a.LoggingOperation(values...)
	responseBody := ResponseToJSON(response.Body.String())
	assert.Equal(t, 201, response.Code, "EXPECTED 201")
	assert.Equal(t, message, responseBody[internal.KeyResponseMessage].(interface{}), "EXPECTED "+message)
	assert.Equal(t, "ANA M.", responseBody[internal.KeyResponseData].(map[string]interface{})[internal.KeyValues].(map[string]interface{})["name"], "EXPECTED ANA M.")
}

func TestDeletePerson(t *testing.T) {
	payload := []byte(`{
						"id" : "60de364abefb1fb4a19d8bb7"
						}`)
	request, _ := http.NewRequest(internal.HTTP_POST, internal.URLDeletingOne, bytes.NewBuffer(payload))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)
	message := func() string {
		if response.Code == 500 {
			return internal.MsgResponseNoData
		} else {
			return internal.MsgResponseUpdatingOne
		}
	}()
	values := []interface{}{internal.KeyType, internal.TEST, internal.KeyURL, internal.URLDeletingOne, internal.KeyMessage, message}
	a.LoggingOperation(values...)
	responseBody := ResponseToJSON(response.Body.String())
	assert.Equal(t, response.Code, response.Code, "EXPECTED "+strconv.Itoa(response.Code))
	assert.Equal(t, message, responseBody[internal.KeyResponseMessage].(interface{}), "EXPECTED "+message)
}
