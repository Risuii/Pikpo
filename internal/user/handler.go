package user

import (
	"encoding/json"
	"net/http"
	"pikpo2/helpers/response"
	"pikpo2/models"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UseCase UserUseCase
}

func NewUserHandler(router *mux.Router, usecase UserUseCase) {
	handler := &UserHandler{
		UseCase: usecase,
	}

	router.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	router.HandleFunc("/profile/{id}", handler.ReadOne).Methods(http.MethodGet)
	router.HandleFunc("/update/{id}", handler.Update).Methods(http.MethodPatch)

}

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var userInput models.User

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		res = response.Error(response.StatusUnprocessableEntity, err)
		res.JSON(w)
		return
	}

	res = handler.UseCase.Register(ctx, userInput)

	res.JSON(w)
}

func (handler *UserHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	ctx := r.Context()
	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 64)

	res = handler.UseCase.ReadOne(ctx, id)

	res.JSON(w)
}

func (handler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var userInput models.User
	ctx := r.Context()

	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 64)

	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		res = response.Error(response.StatusUnprocessableEntity, err)
		res.JSON(w)
		return
	}

	userID := id

	res = handler.UseCase.Update(ctx, id, userID, userInput)

	res.JSON(w)
}
