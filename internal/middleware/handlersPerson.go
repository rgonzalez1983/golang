package middleware

import (
	"encoding/json"
	"go_project/internal/entity"
	"net/http"
	"time"
)

// CRUD

func (a *App) CreatePerson(w http.ResponseWriter, r *http.Request) {
	a.ValidateRequest(w, r, "CreatePerson")
	//request := &entity.Person{}
	data, err, statusCode := a.PersonRepository.CreatePerson(a.FormatRequestPayload(w, r))
	message := "PERSONA INSERTADA CORRECTAMENTE"
	a.FinalResponse(w, message, data, err, statusCode)
}

func (a *App) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	a.ValidateRequest(w, r, "UpdatePerson")
	//request := &entity.Person{}
	data, err, statusCode := a.PersonRepository.UpdatePerson(a.FormatRequestPayload(w, r))
	message := "PERSONA EDITADA CORRECTAMENTE"
	a.FinalResponse(w, message, data, err, statusCode)
}

func (a *App) ListPersons(w http.ResponseWriter, r *http.Request) {
	data, _, statusCode := a.PersonRepository.ListPersons()
	respondWithJSON(w, http.StatusCreated, entity.JsonResponse{
		StatusCode: statusCode,
		Data:       data,
	})
}

func (a *App) GetPerson(w http.ResponseWriter, r *http.Request) {
	a.ValidateRequest(w, r, "GetPerson")
	data, err, statusCode := a.PersonRepository.GetPerson(a.FormatRequestPayload(w, r))
	a.FinalResponse(w, "", data, err, statusCode)
}

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
