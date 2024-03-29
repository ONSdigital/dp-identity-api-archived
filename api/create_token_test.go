package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/api/apitest"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/audit/auditortest"
	"github.com/ONSdigital/go-ns/common"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	authenticateURL = "http://localhost:23800/identity"
)

var (
	testAuthReq = NewTokenRequest{
		Email:    "666@testuser.com",
		Password: "D4mi3n",
	}

	testIdentity = &schema.Identity{
		Name:              "John Paul Jones",
		Email:             "blackdog@ons.gov.uk",
		Password:          "foo",
		UserType:          "bar",
		TemporaryPassword: false,
		Migrated:          false,
		Deleted:           false,
	}

	expectedParams = common.Params{"email": "666@testuser.com"}
)

func TestAPI_AuthenticateEmptyRequestBody(t *testing.T) {
	Convey("should return expected error status if request body is empty", t, func() {
		a := auditortest.New()
		s := &apitest.IdentityServiceMock{}

		r := httptest.NewRequest(http.MethodPost, authenticateURL, nil)
		w := httptest.NewRecorder()
		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.CreateTokenHandler(w, r)

		assertErrorResponse(w.Code, http.StatusBadRequest, w.Body.String(), ErrRequestBodyNil.Error())
		a.AssertRecordCalls()
		So(s.VerifyPasswordCalls(), ShouldHaveLength, 0)
	})
}

func TestAPI_AuthenticationUnsuccessful(t *testing.T) {
	Convey("should return expected error status createToken unsuccessful", t, func() {
		a := auditortest.New()
		s := &apitest.IdentityServiceMock{
			VerifyPasswordFunc: func(ctx context.Context, id string, password string) (*schema.Identity, error) {
				return nil, identity.ErrAuthenticateFailed
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

		authAPI.CreateTokenHandler(w, r)

		assertErrorResponse(w.Code, http.StatusForbidden, w.Body.String(), identity.ErrAuthenticateFailed.Error())

		a.AssertRecordCalls(
			auditortest.Expected{Action: createToken, Result: audit.Attempted, Params: expectedParams},
			auditortest.Expected{Action: createToken, Result: audit.Unsuccessful, Params: expectedParams},
		)
		So(s.VerifyPasswordCalls(), ShouldHaveLength, 1)
		So(s.VerifyPasswordCalls()[0].Email, ShouldEqual, testAuthReq.Email)
		So(s.VerifyPasswordCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticationHandlerIdentityServiceError(t *testing.T) {
	Convey("should return 403 status status if authentication is unsuccessful", t, func() {
		a := auditortest.New()
		s := &apitest.IdentityServiceMock{
			VerifyPasswordFunc: func(ctx context.Context, id string, password string) (*schema.Identity, error) {
				return nil, identity.ErrAuthenticateFailed
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

		authAPI.CreateTokenHandler(w, r)

		assertErrorResponse(w.Code, http.StatusForbidden, w.Body.String(), identity.ErrAuthenticateFailed.Error())
		a.AssertRecordCalls(
			auditortest.Expected{Action: createToken, Result: audit.Attempted, Params: expectedParams},
			auditortest.Expected{Action: createToken, Result: audit.Unsuccessful, Params: expectedParams},
		)
		So(s.VerifyPasswordCalls(), ShouldHaveLength, 1)
		So(s.VerifyPasswordCalls()[0].Email, ShouldEqual, testAuthReq.Email)
		So(s.VerifyPasswordCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticationHandlerUserNotFound(t *testing.T) {
	Convey("should return 404 status status if user not found", t, func() {
		a := auditortest.New()
		s := &apitest.IdentityServiceMock{
			VerifyPasswordFunc: func(ctx context.Context, id string, password string) (*schema.Identity, error) {
				return nil, identity.ErrIdentityNotFound
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

		authAPI.CreateTokenHandler(w, r)

		assertErrorResponse(w.Code, http.StatusNotFound, w.Body.String(), identity.ErrIdentityNotFound.Error())
		a.AssertRecordCalls(
			auditortest.Expected{Action: createToken, Result: audit.Attempted, Params: expectedParams},
			auditortest.Expected{Action: createToken, Result: audit.Unsuccessful, Params: expectedParams},
		)
		So(s.VerifyPasswordCalls(), ShouldHaveLength, 1)
		So(s.VerifyPasswordCalls()[0].Email, ShouldEqual, testAuthReq.Email)
		So(s.VerifyPasswordCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticationHandlerSuccess(t *testing.T) {
	Convey("should return 200 status status if authentication successful", t, func() {
		a := auditortest.New()
		s := &apitest.IdentityServiceMock{
			VerifyPasswordFunc: func(ctx context.Context, id string, password string) (*schema.Identity, error) {
				return testIdentity, nil
			},
		}

		t := &apitest.TokenServiceMock{
			NewTokenFunc: func(ctx context.Context, identity schema.Identity) (*schema.Token, time.Duration, error) {
				return &schema.Token{ID: "666"}, time.Minute * 15, nil
			},
		}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		authAPI := API{
			auditor:         a,
			IdentityService: s,
			Tokens:          t,
		}

		authAPI.CreateTokenHandler(w, r)

		So(w.Code, ShouldEqual, http.StatusOK)

		var authTkn AuthToken
		err = json.Unmarshal(w.Body.Bytes(), &t)
		So(err, ShouldBeNil)
		So(authTkn.Token, ShouldNotBeNil)

		a.AssertRecordCalls(
			auditortest.Expected{Action: createToken, Result: audit.Attempted, Params: expectedParams},
			auditortest.Expected{Action: createToken, Result: audit.Successful, Params: expectedParams},
		)
		So(s.VerifyPasswordCalls(), ShouldHaveLength, 1)
		So(s.VerifyPasswordCalls()[0].Email, ShouldEqual, testAuthReq.Email)
		So(s.VerifyPasswordCalls()[0].Password, ShouldEqual, testAuthReq.Password)
		So(t.NewTokenCalls(), ShouldHaveLength, 1)
		So(t.NewTokenCalls()[0].Identity, ShouldResemble, *testIdentity)
	})
}

func TestAPI_AuthenticateAuditAttemptedError(t *testing.T) {
	Convey("should return expected error status if audit action attempted errors", t, func() {
		a := auditortest.NewErroring(createToken, audit.Attempted)
		s := &apitest.IdentityServiceMock{}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()
		authAPI := API{
			auditor:         a,
			IdentityService: s,
		}

		authAPI.CreateTokenHandler(w, r)

		assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())
		a.AssertRecordCalls(auditortest.Expected{
			Action: createToken,
			Result: audit.Attempted,
			Params: expectedParams,
		})
		So(s.VerifyPasswordCalls(), ShouldHaveLength, 0)
	})
}

func TestAPI_AuthenticationUnsuccessfulAuditUnsuccessfulError(t *testing.T) {
	Convey("should return expected error status if audit action unsuccessful returns an error", t, func() {
		a := auditortest.NewErroring(createToken, audit.Unsuccessful)
		s := &apitest.IdentityServiceMock{
			VerifyPasswordFunc: func(ctx context.Context, id string, password string) (*schema.Identity, error) {
				return nil, identity.ErrAuthenticateFailed
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

		authAPI.CreateTokenHandler(w, r)

		assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())

		a.AssertRecordCalls(
			auditortest.Expected{Action: createToken, Result: audit.Attempted, Params: expectedParams},
			auditortest.Expected{Action: createToken, Result: audit.Unsuccessful, Params: expectedParams},
		)
		So(s.VerifyPasswordCalls(), ShouldHaveLength, 1)
		So(s.VerifyPasswordCalls()[0].Email, ShouldEqual, testAuthReq.Email)
		So(s.VerifyPasswordCalls()[0].Password, ShouldEqual, testAuthReq.Password)
	})
}

func TestAPI_AuthenticationUnsuccessfulAuditSuccessfulError(t *testing.T) {
	Convey("should return expected error status if audit action successful returns an error", t, func() {
		a := auditortest.NewErroring(createToken, audit.Successful)
		s := &apitest.IdentityServiceMock{
			VerifyPasswordFunc: func(ctx context.Context, id string, password string) (*schema.Identity, error) {
				return testIdentity, nil
			},
		}
		t := &apitest.TokenServiceMock{
			NewTokenFunc: func(ctx context.Context, identity schema.Identity) (*schema.Token, time.Duration, error) {
				return &schema.Token{ID: "666"}, time.Minute * 15, nil
			},
		}

		b, err := json.Marshal(testAuthReq)
		So(err, ShouldBeNil)

		r := httptest.NewRequest(http.MethodPost, authenticateURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		authAPI := API{
			auditor:         a,
			IdentityService: s,
			Tokens:          t,
		}

		authAPI.CreateTokenHandler(w, r)

		assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())

		a.AssertRecordCalls(
			auditortest.Expected{Action: createToken, Result: audit.Attempted, Params: expectedParams},
			auditortest.Expected{Action: createToken, Result: audit.Successful, Params: expectedParams},
		)
		So(s.VerifyPasswordCalls(), ShouldHaveLength, 1)
		So(s.VerifyPasswordCalls()[0].Email, ShouldEqual, testAuthReq.Email)
		So(s.VerifyPasswordCalls()[0].Password, ShouldEqual, testAuthReq.Password)
		So(t.NewTokenCalls(), ShouldHaveLength, 1)
		So(t.NewTokenCalls()[0].Identity, ShouldResemble, *testIdentity)
	})
}
