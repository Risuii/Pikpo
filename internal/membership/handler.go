package membership

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MembershipHandler struct {
	UseCase MembershipUseCase
}

func NewMembershipHanlder(router *mux.Router, UseCase MembershipUseCase) {
	handler := &MembershipHandler{
		UseCase: UseCase,
	}

	router.HandleFunc("/membership/{id}/{NumberTier}", handler.Membership).Methods(http.MethodPost)
}

func (handler *MembershipHandler) Membership(w http.ResponseWriter, r *http.Request) {
	// var res response.Response
	ctx := r.Context()

	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Println(err)
	}
	Tier, err := strconv.ParseInt(params["NumberTier"], 10, 64)
	if err != nil {
		log.Println(err)
	}

	_ = handler.UseCase.Membership(ctx, id, Tier)

	// res.JSON(w)
}
