package main

import (
	"github.com/go-chi/chi"
	"github.com/wutchzone/api-response"
	"github.com/wutchzone/auth-service/pkg/accountdb"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func handleRegisterRoute(w http.ResponseWriter, r *http.Request) {
	service := accountdb.NewService(10)
	ServiceDB.SaveAccount(service)
	SessionDB.SetRecord(service.Name(), "service", 0)

	response.CreateResponse(w, response.ResponseOK, service)
}

func handleGetAllServices(w http.ResponseWriter, r *http.Request) {
	cursor := ServiceDB.GetAll()

	if result, err := decodeServices(cursor); err != nil {
		w.WriteHeader(http.StatusNotFound)
		response.CreateResponse(w, response.ResponseError, result)
	} else {
		response.CreateResponse(w, response.ResponseOK, result)
	}
}

func handleGetOneService(w http.ResponseWriter, r *http.Request) {
	dr := ServiceDB.GetAccount(chi.URLParam(r, "id"))

	service := &accountdb.Service{}
	if err := decodeService(dr, service); err != nil {
		w.WriteHeader(http.StatusNotFound)
		response.CreateResponse(w, response.ResponseError, []int{})
	} else {
		response.CreateResponse(w, response.ResponseOK, service)
	}
}

func handleDeleteOneService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := ServiceDB.DeleteAccount(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	SessionDB.RemoveRecord(id)
}

// Decodes multiple services into array
func decodeServices(cursor *mongo.Cursor) ([]accountdb.Service, error) {
	var rslt []accountdb.Service
	decoder := &accountdb.Service{}

	i := 0
	for cursor.Next(nil) {
		err := cursor.Decode(decoder)
		if err != nil {
			return nil, err
		}
		rslt = append(rslt, *decoder)
		i++
	}

	return rslt, nil
}

// Decodes one service
func decodeService(rslt *mongo.SingleResult, decoder *accountdb.Service) error {
	return rslt.Decode(decoder)
}
