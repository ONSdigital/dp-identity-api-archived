package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/store"
	"github.com/ONSdigital/dp-identity-api/store/storetest"
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

		Convey("when createIdentity is called", func() {
			mockStore := &storetest.StorerMock{
				CreateIdentityFunc: func(identity *models.Identity) error {
					return errors.New("expected")
				},
			}

			identityAPI := &IdentityAPI{
				auditor:   auditMock,
				dataStore: store.DataStore{Backend: mockStore},
			}

			r := httptest.NewRequest("POST", "http://localhost:23800/identity", nil)
			w := httptest.NewRecorder()
			identityAPI.CreateIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), auditortest.ErrAudit.Error())
			})

			Convey("and no identity is created", func() {
				So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 0)
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

		Convey("when createIdentity is called", func() {
			mockStore := &storetest.StorerMock{}

			identityAPI := &IdentityAPI{
				auditor:   auditMock,
				dataStore: store.DataStore{Backend: mockStore},
			}

			b, err := json.Marshal([]int{1, 2, 3})
			So(err, ShouldBeNil)

			r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))
			w := httptest.NewRecorder()
			identityAPI.CreateIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusInternalServerError, w.Body.String(), ErrFailedToUnmarshalRequestBody.Error())
			})

			Convey("and no identity is created", func() {
				So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 0)
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

		Convey("when createIdentity is called", func() {
			mockStore := &storetest.StorerMock{
				CreateIdentityFunc: func(identity *models.Identity) error {
					return nil
				},
			}

			identityAPI := &IdentityAPI{
				auditor:   auditMock,
				dataStore: store.DataStore{Backend: mockStore},
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
				So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 1)
				So(mockStore.CreateIdentityCalls()[0].Identity, ShouldResemble, newIdentity)
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

		Convey("when createIdentity is called", func() {
			mockStore := &storetest.StorerMock{
				CreateIdentityFunc: func(identity *models.Identity) error {
					return nil
				},
			}

			identityAPI := &IdentityAPI{
				auditor:   auditMock,
				dataStore: store.DataStore{Backend: mockStore},
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
				So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 1)
				So(mockStore.CreateIdentityCalls()[0].Identity, ShouldResemble, newIdentity)
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
