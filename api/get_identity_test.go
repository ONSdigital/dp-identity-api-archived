package api

import (
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/audit/auditortest"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	defaultUser = &identity.Model{
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
		auditMock := auditortest.NewErroring(getIdentityURL, audit.Attempted)
		serviceMock := &IdentityServiceMock{}

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

			Convey("and no identity is retrieved", func() {
				So(serviceMock.GetCalls(), ShouldHaveLength, 0)
			})

			Convey("and an unsuccessful audit event is recorded", func() {
				auditMock.AssertRecordCalls(auditortest.Expected{Action: getIdentityURL, Result: audit.Attempted, Params: nil})
			})
		})
	})

}


func TestIdentityAPI_GetIdentityError(t *testing.T) {
	Convey("given getIdentity returns an error", t, func() {
		auditMock := auditortest.New()
		serviceMock := &IdentityServiceMock{}

		Convey("when getIdentity is called without a token", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			r := httptest.NewRequest("GET", getIdentityURL, nil)
			w := httptest.NewRecorder()
			identityAPI.GetIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusNotFound, w.Body.String(), ErrNoTokenProvided.Error())
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