package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-identity-api/identity"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WriteEntitySuccess(t *testing.T) {
	Convey("should write expected values to http response", t, func() {
		w := httptest.NewRecorder()
		entity := IdentityCreated{
			ID:  ID,
			URI: createIdentityURL,
		}

		createIdentityResponse.writeEntity(context.Background(), w, entity, http.StatusCreated)

		So(w.Code, ShouldEqual, http.StatusCreated)
		So(w.Header().Get(headerContentType), ShouldEqual, mimeTypeJSON)

		var body IdentityCreated
		err := json.Unmarshal(w.Body.Bytes(), &body)
		So(err, ShouldBeNil)
		So(body, ShouldResemble, entity)
	})
}

func Test_WriteErrorResolveSuccessful(t *testing.T) {
	Convey("should write expected error status and message to http response", t, func() {
		w := httptest.NewRecorder()

		createIdentityResponse.writeError(context.Background(), w, identity.ErrNameValidation)

		So(w.Header().Get(headerContentType), ShouldEqual, "text/plain; charset=utf-8")
		assertErrorResponse(w.Code, http.StatusBadRequest, w.Body.String(), identity.ErrNameValidation.Error())
	})
}

func Test_WriteErrorResolveUnsuccessful(t *testing.T) {
	Convey("should write expected error status and message to http response", t, func() {
		w := httptest.NewRecorder()
		expectedErr := errors.New("wibble")

		createIdentityResponse.writeError(context.Background(), w, expectedErr)

		So(w.Header().Get(headerContentType), ShouldEqual, "text/plain; charset=utf-8")
		assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), expectedErr.Error())
	})
}

func Test_WriteErrorMarshalError(t *testing.T) {
	Convey("should write expected error status and message to http response", t, func() {
		w := httptest.NewRecorder()

		// pass a chan as the entity to cause a json.Marshal error ;-)
		c := make(chan int, 1)
		createIdentityResponse.writeEntity(context.Background(), w, c, http.StatusOK)

		So(w.Header().Get(headerContentType), ShouldEqual, "text/plain; charset=utf-8")
		assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())
	})
}