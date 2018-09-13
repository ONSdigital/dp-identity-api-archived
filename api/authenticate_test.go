package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/audit/auditortest"
	"github.com/ONSdigital/go-ns/common"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	authenticateURL = "http://localhost:23800/identity"
)

var (
	testAuthReq = AuthenticateRequest{
		ID:       "666",
		Password: "D4mi3n",
	}
)

func TestAPI_AuthenticateEmptyRequestBody(t *testing.T) {
	Convey("should return expected error status if request body is empty", t, func() {
		a := auditortest.New()
		s := &IdentityServiceMock{}

		r := httptest.NewRequest(http.MethodPost, authenticateURL, nil)
		w := httptest.NewRecorder()
		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.AuthenticationHandler(w, r)

		assertErrorResponse(w.Code, http.StatusBadRequest, w.Body.String(), ErrRequestBodyNil.Error())
		a.AssertRecordCalls()
		So(s.AuthenticateCalls(), ShouldHaveLength, 0)
	})
}

func TestAPI_AuthenticationUnsuccessful(t *testing.T) {
	Convey("should return expected error status authenticate unsuccessful", t, func() {
		a := auditortest.New()
		s := &IdentityServiceMock{
			AuthenticateFunc: func(ctx context.Context, id string, password string) error {
				return identity.ErrAuthenticateFailed
			},
		}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.AuthenticationHandler(w, r)

		assertErrorResponse(w.Code, http.StatusForbidden, w.Body.String(), identity.ErrAuthenticateFailed.Error())

		a.AssertRecordCalls(
			auditortest.Expected{Action: authenticateAction, Result: audit.Attempted, Params: common.Params{"id": "666"}},
			auditortest.Expected{Action: authenticateAction, Result: audit.Unsuccessful, Params: common.Params{"id": "666"}},
		)
		So(s.AuthenticateCalls(), ShouldHaveLength, 1)
		So(s.AuthenticateCalls()[0].ID, ShouldEqual, testAuthReq.ID)
		So(s.AuthenticateCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticationHandlerIdentityServiceError(t *testing.T) {
	Convey("should return 403 status status if authentication is unsuccessful", t, func() {
		a := auditortest.New()
		s := &IdentityServiceMock{
			AuthenticateFunc: func(ctx context.Context, id string, password string) error {
				return identity.ErrAuthenticateFailed
			},
		}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.AuthenticationHandler(w, r)

		assertErrorResponse(w.Code, http.StatusForbidden, w.Body.String(), identity.ErrAuthenticateFailed.Error())
		a.AssertRecordCalls(
			auditortest.Expected{Action: authenticateAction, Result: audit.Attempted, Params: common.Params{"id": "666"}},
			auditortest.Expected{Action: authenticateAction, Result: audit.Unsuccessful, Params: common.Params{"id": "666"}},
		)
		So(s.AuthenticateCalls(), ShouldHaveLength, 1)
		So(s.AuthenticateCalls()[0].ID, ShouldEqual, testAuthReq.ID)
		So(s.AuthenticateCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticationHandlerUserNotFound(t *testing.T) {
	Convey("should return 404 status status if user not found", t, func() {
		a := auditortest.New()
		s := &IdentityServiceMock{
			AuthenticateFunc: func(ctx context.Context, id string, password string) error {
				return identity.ErrUserNotFound
			},
		}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.AuthenticationHandler(w, r)

		assertErrorResponse(w.Code, http.StatusNotFound, w.Body.String(), identity.ErrUserNotFound.Error())
		a.AssertRecordCalls(
			auditortest.Expected{Action: authenticateAction, Result: audit.Attempted, Params: common.Params{"id": "666"}},
			auditortest.Expected{Action: authenticateAction, Result: audit.Unsuccessful, Params: common.Params{"id": "666"}},
		)
		So(s.AuthenticateCalls(), ShouldHaveLength, 1)
		So(s.AuthenticateCalls()[0].ID, ShouldEqual, testAuthReq.ID)
		So(s.AuthenticateCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticationHandlerSuccess(t *testing.T) {
	Convey("should return 200 status status if authentication successful", t, func() {
		a := auditortest.New()
		s := &IdentityServiceMock{
			AuthenticateFunc: func(ctx context.Context, id string, password string) error {
				return nil
			},
		}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.AuthenticationHandler(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)

		var t AuthToken
		err = json.Unmarshal(w.Body.Bytes(), &t)
		So(err, ShouldBeNil)
		So(t.Token, ShouldNotBeNil)

		a.AssertRecordCalls(
			auditortest.Expected{Action: authenticateAction, Result: audit.Attempted, Params: common.Params{"id": "666"}},
			auditortest.Expected{Action: authenticateAction, Result: audit.Successful, Params: common.Params{"id": "666"}},
		)
		So(s.AuthenticateCalls(), ShouldHaveLength, 1)
		So(s.AuthenticateCalls()[0].ID, ShouldEqual, testAuthReq.ID)
		So(s.AuthenticateCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticateAuditAttemptedError(t *testing.T) {
	Convey("should return expected error status if audit action attempted errors", t, func() {
		a := auditortest.NewErroring(authenticateAction, audit.Attempted)
		s := &IdentityServiceMock{}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()
		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.AuthenticationHandler(w, r)

		assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())
		a.AssertRecordCalls(auditortest.Expected{
			Action: authenticateAction,
			Result: audit.Attempted,
			Params: common.Params{"id": "666"},
		})
		So(s.AuthenticateCalls(), ShouldHaveLength, 0)
	})
}

func TestAPI_AuthenticationUnsuccessfulAuditUnsuccessfulError(t *testing.T) {
	Convey("should return expected error status if audit action unsuccessful returns an error", t, func() {
		a := auditortest.NewErroring(authenticateAction, audit.Unsuccessful)
		s := &IdentityServiceMock{
			AuthenticateFunc: func(ctx context.Context, id string, password string) error {
				return identity.ErrAuthenticateFailed
			},
		}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.AuthenticationHandler(w, r)

		assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())

		a.AssertRecordCalls(
			auditortest.Expected{Action: authenticateAction, Result: audit.Attempted, Params: common.Params{"id": "666"}},
			auditortest.Expected{Action: authenticateAction, Result: audit.Unsuccessful, Params: common.Params{"id": "666"}},
		)
		So(s.AuthenticateCalls(), ShouldHaveLength, 1)
		So(s.AuthenticateCalls()[0].ID, ShouldEqual, testAuthReq.ID)
		So(s.AuthenticateCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticationUnsuccessfulAuditSuccessfulError(t *testing.T) {
	Convey("should return expected error status if audit action successful returns an error", t, func() {
		a := auditortest.NewErroring(authenticateAction, audit.Successful)
		s := &IdentityServiceMock{
			AuthenticateFunc: func(ctx context.Context, id string, password string) error {
				return nil
			},
		}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.AuthenticationHandler(w, r)

		assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())

		a.AssertRecordCalls(
			auditortest.Expected{Action: authenticateAction, Result: audit.Attempted, Params: common.Params{"id": "666"}},
			auditortest.Expected{Action: authenticateAction, Result: audit.Successful, Params: common.Params{"id": "666"}},
		)
		So(s.AuthenticateCalls(), ShouldHaveLength, 1)
		So(s.AuthenticateCalls()[0].ID, ShouldEqual, testAuthReq.ID)
		So(s.AuthenticateCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}
