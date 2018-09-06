package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/identity"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_writeJSONBodySuccess(t *testing.T) {
	Convey("should write expected values to http response", t, func() {
		w := httptest.NewRecorder()

		ident := identity.Model{
			Name:     "Solid Snake",
			Email:    "snake@mgs.com",
			UserType: "FoxHound",
			Password: "M3t4l G34r S0L1D",
		}

		writeJSONBody(context.Background(), w, ident, http.StatusCreated)

		So(w.Code, ShouldEqual, http.StatusCreated)
		So(w.Header().Get(headerContentType), ShouldEqual, mimeTypeJSON)

		var body identity.Model
		err := json.Unmarshal(w.Body.Bytes(), &body)
		So(err, ShouldBeNil)
		So(body, ShouldResemble, ident)
	})
}

func Test_writeJSONBodyMarshalError(t *testing.T) {
	Convey("should write http 500 response if marshal returns an error", t, func() {
		w := httptest.NewRecorder()

		// pass in a channel to cause json.Marshal to return an error
		input := make(chan int, 1)

		writeJSONBody(context.Background(), w, input, http.StatusCreated)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(w.Header().Get(headerContentType), ShouldEqual, "text/plain; charset=utf-8")
	})
}
