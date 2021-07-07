package middleware

import (
	"encoding/json"
	"go_project/internal"
	"go_project/internal/entity"
	"net/http"
	"time"
)

// CRUD

// CreatePerson godoc
// @Summary Create one person
// @Description Create one person
// @Tags CRUD
// @Accept  json
// @Produce  json
// @Param parameters body entity.Person true "Create Person"
// @Success 200 {object} entity.InterfaceAPI
// @Router /create_person [post]
func (a *App) CreatePerson(w http.ResponseWriter, r *http.Request) {
	a.ValidateRequest(w, r, internal.FuncCreatePerson)
	data, err, statusCode := a.PersonRepository.CreatePerson(a.FormatRequestPayload(w, r))
	indicator := a.IndicatorType(statusCode)
	message := func() string {
		if indicator == internal.ERROR {
			return err.Error()
		} else {
			return internal.MsgResponseCreatingOne
		}
	}()
	values := []interface{}{internal.KeyType, indicator, internal.KeyURL, internal.URLCreatingOne, internal.KeyMessage, message}
	a.LoggingOperation(values...)
	a.FinalResponse(w, message, data, err, statusCode)
}

// UpdatePerson godoc
// @Summary Update one person
// @Description Update of one person
// @Tags CRUD
// @Accept  json
// @Produce  json
// @Param parameters body entity.InterfaceAPI true "Update Person"
// @Success 200 {object} entity.InterfaceAPI
// @Router /update_person [post]
func (a *App) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	a.ValidateRequest(w, r, internal.FuncUpdatePerson)
	data, err, statusCode := a.PersonRepository.UpdatePerson(a.FormatRequestPayload(w, r))
	indicator := a.IndicatorType(statusCode)
	message := func() string {
		if indicator == internal.ERROR {
			return err.Error()
		} else {
			return internal.MsgResponseUpdatingOne
		}
	}()
	values := []interface{}{internal.KeyType, indicator, internal.KeyURL, internal.URLUpdatingOne, internal.KeyMessage, message}
	a.LoggingOperation(values...)
	a.FinalResponse(w, message, data, err, statusCode)
}

// ListPerson godoc
// @Summary Get details of all persons
// @Description Get details of all persons
// @Tags CRUD
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.InterfaceAPI
// @Router /list_persons [post]
func (a *App) ListPersons(w http.ResponseWriter, r *http.Request) {
	a.ValidateRequest(w, r, internal.FuncListPersons)
	data, _, statusCode := a.PersonRepository.ListPersons()
	indicator := a.IndicatorType(statusCode)
	values := []interface{}{internal.KeyType, indicator, internal.KeyURL, internal.URLListingAll, internal.KeyMessage, internal.MsgResponseListingAll}
	a.LoggingOperation(values...)
	a.FinalResponse(w, internal.ValueEmpty, data, nil, statusCode)
}

// GetPerson godoc
// @Summary Get details of one person
// @Description Get details of one person
// @Tags CRUD
// @Accept  json
// @Produce  json
// @Param parameters body entity.InterfaceAPI true "Get Person"
// @Success 200 {object} entity.InterfaceAPI
// @Router /get_person [post]
func (a *App) GetPerson(w http.ResponseWriter, r *http.Request) {
	a.ValidateRequest(w, r, internal.FuncGetPerson)
	data, err, statusCode := a.PersonRepository.GetPerson(a.FormatRequestPayload(w, r))
	indicator := a.IndicatorType(statusCode)
	message := func() string {
		if indicator == internal.ERROR {
			return err.Error()
		} else {
			return internal.MsgResponseGettingOne
		}
	}()
	values := []interface{}{internal.KeyType, indicator, internal.KeyURL, internal.URLGettingOne, internal.KeyMessage, message}
	a.LoggingOperation(values...)
	a.FinalResponse(w, message, data, err, statusCode)
}

// DeletePerson godoc
// @Summary Delete one person
// @Description Delete of one person
// @Tags CRUD
// @Accept  json
// @Produce  json
// @Param parameters body entity.InterfaceAPI true "Delete Person"
// @Success 200 {object} entity.InterfaceAPI
// @Router /delete_person [post]
func (a *App) DeletePerson(w http.ResponseWriter, r *http.Request) {
	a.ValidateRequest(w, r, internal.FuncDeletePerson)
	data, err, statusCode := a.PersonRepository.DeletePerson(a.FormatRequestPayload(w, r))
	indicator := a.IndicatorType(statusCode)
	message := func() string {
		if indicator == internal.ERROR {
			return err.Error()
		} else {
			return internal.MsgResponseDeletingOne
		}
	}()
	values := []interface{}{internal.KeyType, indicator, internal.KeyURL, internal.URLDeletingOne, internal.KeyMessage, message}
	a.LoggingOperation(values...)
	a.FinalResponse(w, message, data, err, statusCode)
}

// FUNCIONES EXTRAS

func (a *App) ValidateRequest(w http.ResponseWriter, r *http.Request, function string) {
	if (*r).Method == internal.OPTIONS {
		respondWithJSON(w, http.StatusInternalServerError, entity.JsonResponse{
			Message:    internal.MsgResponseServerError,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	defer func(begin time.Time) {
		values := []interface{}{internal.KeyMethod, function, internal.KeyTook, time.Since(begin)}
		//_ = a.Logg.Log(values)
		a.LoggingOperation(values...)
	}(time.Now())
}

func (a *App) FormatRequestPayload(w http.ResponseWriter, r *http.Request) (request *interface{}) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondWithJSON(w, http.StatusBadRequest, entity.JsonResponse{
			Message:    internal.MsgResponseInvalidRequest,
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	defer r.Body.Close()
	return request
}

func (a *App) FinalResponse(w http.ResponseWriter, message string, data interface{}, err error, statusCode int) {
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, entity.JsonResponse{
			Message:    err.Error(),
			StatusCode: statusCode,
			Data:       data,
		})
		return
	}
	respondWithJSON(w, http.StatusCreated, entity.JsonResponse{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	})
}

func (a *App) IndicatorType(statusCode int) string {
	return func() string {
		if statusCode == 500 {
			return internal.ERROR
		} else {
			return internal.SUCCESS
		}
	}()
}
