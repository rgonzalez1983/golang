package middleware

import (
	"encoding/json"
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
	a.ValidateRequest(w, r, "CreatePerson")
	//request := &entity.Person{}
	data, err, statusCode := a.PersonRepository.CreatePerson(a.FormatRequestPayload(w, r))
	message := "PERSONA INSERTADA CORRECTAMENTE"
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
	a.ValidateRequest(w, r, "UpdatePerson")
	//request := &entity.Person{}
	data, err, statusCode := a.PersonRepository.UpdatePerson(a.FormatRequestPayload(w, r))
	message := "PERSONA EDITADA CORRECTAMENTE"
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
	data, _, statusCode := a.PersonRepository.ListPersons()
	respondWithJSON(w, http.StatusCreated, entity.JsonResponse{
		StatusCode: statusCode,
		Data:       data,
	})
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
	a.ValidateRequest(w, r, "GetPerson")
	data, err, statusCode := a.PersonRepository.GetPerson(a.FormatRequestPayload(w, r))
	a.FinalResponse(w, "", data, err, statusCode)
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
	a.ValidateRequest(w, r, "DeletePerson")
	//request := &entity.Person{}
	data, err, statusCode := a.PersonRepository.DeletePerson(a.FormatRequestPayload(w, r))
	message := "PERSONA BORRADA CORRECTAMENTE"
	a.FinalResponse(w, message, data, err, statusCode)
}

// FUNCIONES EXTRAS

func (a *App) ValidateRequest(w http.ResponseWriter, r *http.Request, function string) {
	if (*r).Method == "OPTIONS" {
		respondWithJSON(w, http.StatusInternalServerError, entity.JsonResponse{
			Message:    "ERROR EN LA PETICION",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	defer func(begin time.Time) {
		a.Logg.Log("method", function, "took", time.Since(begin))
	}(time.Now())
}

func (a *App) FormatRequestPayload(w http.ResponseWriter, r *http.Request) (request *interface{}) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondWithJSON(w, http.StatusBadRequest, entity.JsonResponse{
			Message:    "PAYLOAD/REQUEST INVALIDO",
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
