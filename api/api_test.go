package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/audit/auditortest"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIdentityAPI_CreateIdentityAuditAttemptedFailed(t *testing.T) {
	Convey("given audit action attempted returns an error", t, func() {
		auditMock := auditortest.NewErroring(createIdentityAction, audit.Attempted)
		serviceMock := &IdentityServiceMock{}

		Convey("when createIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			r := httptest.NewRequest("POST", "http://localhost:23800/identity", nil)
			w := httptest.NewRecorder()
			identityAPI.CreateIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), auditortest.ErrAudit.Error())
			})

			Convey("and no identity is created", func() {
				So(serviceMock.CreateCalls(), ShouldHaveLength, 0)
			})

			Convey("and an unsuccessful audit event is recorded", func() {
				auditMock.AssertRecordCalls(auditortest.Expected{Action: createIdentityAction, Result: audit.Attempted, Params: nil})
			})
		})
	})
}

func TestIdentityAPI_CreateIdentityError(t *testing.T) {
	Convey("given createIdentity returns an error", t, func() {
		auditMock := auditortest.New()
		serviceMock := &IdentityServiceMock{
			CreateFunc: func(ctx context.Context, r *http.Request) error {
				return identity.ErrFailedToUnmarshalRequestBody
			},
		}

		Convey("when createIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			b, err := json.Marshal([]int{1, 2, 3})
			So(err, ShouldBeNil)

			r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))
			w := httptest.NewRecorder()
			identityAPI.CreateIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), identity.ErrFailedToUnmarshalRequestBody.Error())
			})

			Convey("and no identity is created", func() {
				So(serviceMock.CreateCalls(), ShouldHaveLength, 1)
			})

			Convey("and an unsuccessful audit event is recorded", func() {
				auditMock.AssertRecordCalls(
					auditortest.Expected{Action: createIdentityAction, Result: audit.Attempted, Params: nil},
					auditortest.Expected{Action: createIdentityAction, Result: audit.Unsuccessful, Params: nil},
				)
			})
		})
	})
}

func TestIdentityAPI_CreateIdentityAuditSuccessfulError(t *testing.T) {
	Convey("given audit action successful returns an error", t, func() {
		auditMock := auditortest.NewErroring(createIdentityAction, audit.Successful)
		serviceMock := &IdentityServiceMock{
			CreateFunc: func(ctx context.Context, r *http.Request) error {
				return nil
			},
		}

		Convey("when createIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			newIdentity := &models.Identity{Name: "Eleven"}
			b, err := json.Marshal(newIdentity)
			So(err, ShouldBeNil)

			r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))
			w := httptest.NewRecorder()
			identityAPI.CreateIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), auditortest.ErrAudit.Error())
			})

			Convey("and the identity is created", func() {
				So(serviceMock.CreateCalls(), ShouldHaveLength, 1)
			})

			Convey("and attempted and successful audit events are recorded", func() {
				auditMock.AssertRecordCalls(
					auditortest.Expected{Action: createIdentityAction, Result: audit.Attempted, Params: nil},
					auditortest.Expected{Action: createIdentityAction, Result: audit.Successful, Params: nil},
				)
			})
		})
	})
}

func TestIdentityAPI_CreateIdentitySuccess(t *testing.T) {
	Convey("given create identity is successful", t, func() {
		auditMock := auditortest.New()
		serviceMock := &IdentityServiceMock{
			CreateFunc: func(ctx context.Context, r *http.Request) error {
				return nil
			},
		}

		Convey("when createIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			newIdentity := &models.Identity{Name: "Eleven"}
			b, err := json.Marshal(newIdentity)
			So(err, ShouldBeNil)

			r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))
			w := httptest.NewRecorder()
			identityAPI.CreateIdentityHandler(w, r)

			Convey("then a HTTP 200 status is returned", func() {
				assertErrorResponse(w.Code, http.StatusCreated, w.Body.String(), "")
			})

			Convey("and the identity is created", func() {
				So(serviceMock.CreateCalls(), ShouldHaveLength, 1)
			})

			Convey("and attempted and successful audit events are recorded", func() {
				auditMock.AssertRecordCalls(
					auditortest.Expected{Action: createIdentityAction, Result: audit.Attempted, Params: nil},
					auditortest.Expected{Action: createIdentityAction, Result: audit.Successful, Params: nil},
				)
			})
		})
	})
}

func assertErrorResponse(actualStatus int, expectedStatus int, actualBody string, expectedBody string) {
	So(actualStatus, ShouldEqual, expectedStatus)
	So(strings.TrimSpace(actualBody), ShouldEqual, expectedBody)
}
