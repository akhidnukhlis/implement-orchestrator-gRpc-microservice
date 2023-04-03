package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/akhidnukhlis/implement-gRpc-microservice-orchestrator/helpers"
	"github.com/akhidnukhlis/implement-gRpc-microservice-orchestrator/internal/entity"
	"github.com/akhidnukhlis/implement-gRpc-microservice-orchestrator/internal/entity/author_entity"
	"github.com/akhidnukhlis/implement-gRpc-microservice-orchestrator/internal/src/author"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type AuthorHandler struct {
	service author.Service
}

func NewAuthorHandler(service author.Service) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

// RegisterNewAuthor is handler function to Handle author registration
func (ah *AuthorHandler) RegisterNewAuthor(w http.ResponseWriter, r *http.Request) {
	responder := helpers.NewHTTPResponse("registerNewAuthor")
	ctx := r.Context()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responder.ErrorWithStatusCode(w, http.StatusUnprocessableEntity, fmt.Sprint(err))
		return
	}

	var payload author_entity.AuthorRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		responder.ErrorWithStatusCode(w, http.StatusUnprocessableEntity, fmt.Sprint(err))
		return
	}

	err = ah.service.InsertNewAuthor(ctx, &payload)
	if err != nil {
		causer := errors.Cause(err)
		switch causer {
		case entity.ErrAuthorAlreadyExist:
			responder.FieldErrors(w, err, http.StatusNotAcceptable, err.Error())
			return
		default:
			responder.FieldErrors(w, err, http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
	}
	responder.SuccessWithoutData(w, http.StatusCreated, "Success to register new author")
	return
}

// FindUserByAuthorID is handler function to Handle find author
func (ah *AuthorHandler) FindUserByAuthorID(w http.ResponseWriter, r *http.Request) {
	var (
		authorID  = mux.Vars(r)["id"]
		responder = helpers.NewHTTPResponse("registerNewAuthor")
		ctx       = r.Context()
	)

	findAuthor, err := ah.service.FindAuthor(ctx, authorID)
	if err != nil {
		causer := errors.Cause(err)
		switch causer {
		case entity.ErrAuthorNotExist:
			responder.ErrorJSON(w, http.StatusNotFound, "author not found")
			return
		default:
			responder.FailureJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	responder.SuccessJSON(w, findAuthor, http.StatusOK, "Author found")
	return
}
