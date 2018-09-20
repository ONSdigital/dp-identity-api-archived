package api

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/api/apitest"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/audit/auditortest"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/ONSdigital/dp-identity-api/schema"
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
)

const getIdentityURL = "http://localhost:23800/identity"

func TestIdentityAPI_GetIdentityAuditAttemptedFailed(t *testing.T) {
	Convey("given audit action attempted returns an error", t, func() {
		auditMock := auditortest.NewErroring(getIdentityAction, audit.Attempted)
		serviceMock := &apitest.IdentityServiceMock{}

		Convey("when getIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), auditortest.ErrAudit.Error())
			})

			Convey("and no identity is created", func() {
				So(serviceMock.GetCalls(), ShouldHaveLength, 0)
			})

			Convey("and an unsuccessful audit event is recorded", func() {
				auditMock.AssertRecordCalls(auditortest.Expected{Action: getIdentityAction, Result: audit.Attempted, Params: nil})
			})
		})
	})
}

func TestIdentityAPI_GetIdentityError(t *testing.T) {
	Convey("given getIdentity returns an error", t, func() {
		auditMock := auditortest.New()
		serviceMock := &apitest.IdentityServiceMock{}

		Convey("when getIdentity is called without a token", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusUnauthorized, w.Body.String(), ErrNoTokenProvided.Error())
			})

			Convey("and identity service is never called", func() {
				So(serviceMock.GetCalls(), ShouldHaveLength, 0)
			})

			Convey("and an unsuccessful audit event is recorded", func() {
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
		serviceMock := &apitest.IdentityServiceMock{
			GetFunc: func(ctx context.Context, tokenStr string) (*schema.Identity, error) {
				return defaultUser, nil
			},
		}

		Convey("when getIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			r.Header.Set("token", "1234")
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), auditortest.ErrAudit.Error())
			})

			Convey("and the identity is created", func() {
				So(serviceMock.GetCalls(), ShouldHaveLength, 1)
			})

			Convey("and attempted and successful audit events are recorded", func() {
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
		serviceMock := &apitest.IdentityServiceMock{
			GetFunc: func(ctx context.Context, tokenStr string) (*schema.Identity, error) {
				return defaultUser, nil
			},
		}

		Convey("when getIdentity is called", func() {
			identityAPI := &API{
				Host:            "http://localhost:23800",
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			r.Header.Set("token", "1234")
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then a HTTP 200 status is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("and the identity is retrieved", func() {
				So(serviceMock.GetCalls(), ShouldHaveLength, 1)
			})

			Convey("and attempted and successful audit events are recorded", func() {
				auditMock.AssertRecordCalls(
					auditortest.Expected{Action: getIdentityAction, Result: audit.Attempted, Params: nil},
					auditortest.Expected{Action: getIdentityAction, Result: audit.Successful, Params: nil},
				)
			})
		})
	})
}

func TestGetIdentity_IdentityServiceError(t *testing.T) {
	Convey("should return expected error if identityService returns an error", t, func() {
		serviceMock := &apitest.IdentityServiceMock{
			GetFunc: func(ctx context.Context, tokenStr string) (*schema.Identity, error) {
				return nil, errTest
			},
		}

		identityAPI := &API{IdentityService: serviceMock}

		r := httptest.NewRequest("GET", getIdentityURL, nil)
		r.Header.Set("token", "1234")
		i, err := identityAPI.getIdentity(context.Background(), r)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, errTest)
		So(i, ShouldBeNil)
	})
}

func TestGetIdentity_NoTokenProvided(t *testing.T) {
	Convey("should return expected error if no token is provided in the request header", t, func() {
		serviceMock := &apitest.IdentityServiceMock{}
		identityAPI := &API{IdentityService: serviceMock}

		r := httptest.NewRequest("GET", createIdentityURL, nil)
		i, err := identityAPI.getIdentity(context.Background(), r)

		So(err, ShouldEqual, ErrNoTokenProvided)
		So(i, ShouldBeNil)
		So(serviceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}
