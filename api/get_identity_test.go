package api

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/api/apitest"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/audit/auditortest"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	defaultUser = &schema.Identity{
		Name:              "John Paul Jones",
		Email:             "blackdog@ons.gov.uk",
		Password:          "foo",
		UserType:          "bar",
		TemporaryPassword: false,
		Migrated:          false,
		Deleted:           false,
	}

	tkn = &schema.Token{
		ID:         "666",
		IdentityID: "666",
	}

	tokenTTL = time.Minute * 10
)

const getIdentityURL = "http://localhost:23800/identity"

func TestIdentityAPI_GetIdentityAuditAttemptedFailed(t *testing.T) {
	Convey("given audit action attempted returns an error", t, func() {
		auditMock := auditortest.NewErroring(getIdentityAction, audit.Attempted)
		tokens := &apitest.TokenServiceMock{}

		Convey("when getIdentity is called", func() {
			identityAPI := &API{
				auditor: auditMock,
				Tokens:  tokens,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())
				So(tokens.GetIdentityByTokenCalls(), ShouldHaveLength, 0)
				auditMock.AssertRecordCalls(auditortest.Expected{Action: getIdentityAction, Result: audit.Attempted, Params: nil})
			})
		})
	})
}

func TestIdentityAPI_GetIdentityError(t *testing.T) {
	Convey("given getIdentity returns an error", t, func() {
		auditMock := auditortest.New()
		tokens := &apitest.TokenServiceMock{}

		Convey("when getIdentity is called without a token", func() {
			identityAPI := &API{
				auditor: auditMock,
				Tokens:  tokens,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusUnauthorized, w.Body.String(), ErrNoTokenProvided.Error())
				So(tokens.GetIdentityByTokenCalls(), ShouldHaveLength, 0)

				auditMock.AssertRecordCalls(
					auditortest.Expected{Action: getIdentityAction, Result: audit.Attempted, Params: nil},
					auditortest.Expected{Action: getIdentityAction, Result: audit.Unsuccessful, Params: nil},
				)
			})
		})
	})
}

func TestIdentityAPI_GetIdentityAuditSuccessfulError(t *testing.T) {
	Convey("given audit action successful returns an error", t, func() {
		auditMock := auditortest.NewErroring(getIdentityAction, audit.Successful)

		tokensMock := &apitest.TokenServiceMock{
			GetIdentityByTokenFunc: func(ctx context.Context, tokenStr string) (*schema.Identity, time.Duration, error) {
				return defaultUser, tokenTTL, nil
			},
		}

		Convey("when getIdentity is called", func() {
			identityAPI := &API{
				auditor: auditMock,
				Tokens:  tokensMock,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			r.Header.Set("token", "1234")
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrInternalServerError.Error())
				So(tokensMock.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(tokensMock.GetIdentityByTokenCalls()[0].TokenStr, ShouldEqual, "1234")

				auditMock.AssertRecordCalls(
					auditortest.Expected{Action: getIdentityAction, Result: audit.Attempted, Params: nil},
					auditortest.Expected{Action: getIdentityAction, Result: audit.Successful, Params: nil},
				)
			})
		})
	})
}

func TestIdentityAPI_GetIdentitySuccess(t *testing.T) {
	Convey("given create identity is successful", t, func() {
		auditMock := auditortest.New()
		tokensMock := &apitest.TokenServiceMock{
			GetIdentityByTokenFunc: func(ctx context.Context, tokenStr string) (*schema.Identity, time.Duration, error) {
				return defaultUser, tokenTTL, nil
			},
		}

		Convey("when getIdentity is called", func() {
			identityAPI := &API{
				Host:    "http://localhost:23800",
				auditor: auditMock,
				Tokens:  tokensMock,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			r.Header.Set("token", "1234")
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then a HTTP 200 status is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)

				So(tokensMock.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
				So(tokensMock.GetIdentityByTokenCalls()[0].TokenStr, ShouldEqual, "1234")

				auditMock.AssertRecordCalls(
					auditortest.Expected{Action: getIdentityAction, Result: audit.Attempted, Params: nil},
					auditortest.Expected{Action: getIdentityAction, Result: audit.Successful, Params: nil},
				)
			})
		})
	})
}

func TestGetIdentity_IdentityServiceError(t *testing.T) {
	Convey("should return expected error if tokens returns an error", t, func() {
		tokensMock := &apitest.TokenServiceMock{
			GetIdentityByTokenFunc: func(ctx context.Context, tokenStr string) (*schema.Identity, time.Duration, error) {
				return nil, 0, schema.ErrTokenExpired
			},
		}

		identityAPI := &API{Tokens: tokensMock}

		r := httptest.NewRequest("GET", getIdentityURL, nil)
		r.Header.Set("token", "1234")
		i, err := identityAPI.getIdentity(context.Background(), r)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, schema.ErrTokenExpired)
		So(i, ShouldBeNil)
		So(tokensMock.GetIdentityByTokenCalls(), ShouldHaveLength, 1)
		So(tokensMock.GetIdentityByTokenCalls()[0].TokenStr, ShouldEqual, "1234")
	})
}

func TestGetIdentity_NoTokenProvided(t *testing.T) {
	Convey("should return expected error if no token is provided in the request header", t, func() {
		tokensMock := &apitest.TokenServiceMock{}
		identityAPI := &API{Tokens: tokensMock}

		r := httptest.NewRequest("GET", createIdentityURL, nil)
		i, err := identityAPI.getIdentity(context.Background(), r)

		So(err, ShouldEqual, ErrNoTokenProvided)
		So(i, ShouldBeNil)
		So(tokensMock.GetIdentityByTokenCalls(), ShouldHaveLength, 0)
	})
}
