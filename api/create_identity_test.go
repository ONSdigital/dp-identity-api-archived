package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/api/apitest"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/audit/auditortest"
	"github.com/ONSdigital/go-ns/common"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	ID = "666"

	createIdentityURL = "http://localhost:23800/identity"
)

var (
	expectedIdentity = IdentityCreated{
		ID:  ID,
		URI: "http://localhost:23800/identity/" + ID,
	}

	errTest = errors.New("boom!")
)

//IOReaderErroring is an io.Reader for unit testing that returns the specified error when Read() is called.
type IOReaderErroring struct {
	err error
}

func (r *IOReaderErroring) Read(p []byte) (int, error) {
	return 0, r.err
}

func TestIdentityAPI_CreateIdentityAuditAttemptedFailed(t *testing.T) {
	Convey("given audit action attempted returns an error", t, func() {
		auditMock := auditortest.NewErroring(createIdentityAction, audit.Attempted)
		serviceMock := &apitest.IdentityServiceMock{}

		Convey("when createIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			r := httptest.NewRequest("POST", createIdentityURL, nil)
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
		serviceMock := &apitest.IdentityServiceMock{}

		Convey("when createIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			b, err := json.Marshal([]int{1, 2, 3})
			So(err, ShouldBeNil)

			r := httptest.NewRequest("POST", createIdentityURL, bytes.NewReader(b))
			w := httptest.NewRecorder()
			identityAPI.CreateIdentityHandler(w, r)

			Convey("then the expected error response is returned", func() {
				assertErrorResponse(w.Code, http.StatusBadRequest, w.Body.String(), ErrFailedToUnmarshalRequestBody.Error())
			})

			Convey("and identity service is never called", func() {
				So(serviceMock.CreateCalls(), ShouldHaveLength, 0)
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
		serviceMock := &apitest.IdentityServiceMock{
			CreateFunc: func(ctx context.Context, i *identity.Model) (string, error) {
				return ID, nil
			},
		}

		Convey("when createIdentity is called", func() {
			identityAPI := &API{
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			newIdentity := &identity.Model{Name: "Eleven"}
			b, err := json.Marshal(newIdentity)
			So(err, ShouldBeNil)

			r := httptest.NewRequest("POST", createIdentityURL, bytes.NewReader(b))
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
					auditortest.Expected{Action: createIdentityAction, Result: audit.Successful, Params: common.Params{"id": ID}},
				)
			})
		})
	})
}

func TestIdentityAPI_CreateIdentitySuccess(t *testing.T) {
	Convey("given create identity is successful", t, func() {
		auditMock := auditortest.New()
		serviceMock := &apitest.IdentityServiceMock{
			CreateFunc: func(ctx context.Context, i *identity.Model) (string, error) {
				return ID, nil
			},
		}

		Convey("when createIdentity is called", func() {
			identityAPI := &API{
				Host:            "http://localhost:23800",
				auditor:         auditMock,
				IdentityService: serviceMock,
			}

			newIdentity := &identity.Model{Name: "Eleven"}
			b, err := json.Marshal(newIdentity)
			So(err, ShouldBeNil)

			r := httptest.NewRequest("POST", createIdentityURL, bytes.NewReader(b))
			w := httptest.NewRecorder()
			identityAPI.CreateIdentityHandler(w, r)

			Convey("then a HTTP 201 status is returned", func() {
				So(w.Code, ShouldEqual, http.StatusCreated)

				var actual IdentityCreated
				err := json.Unmarshal(w.Body.Bytes(), &actual)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expectedIdentity)

			})

			Convey("and the identity is created", func() {
				So(serviceMock.CreateCalls(), ShouldHaveLength, 1)
			})

			Convey("and attempted and successful audit events are recorded", func() {
				auditMock.AssertRecordCalls(
					auditortest.Expected{Action: createIdentityAction, Result: audit.Attempted, Params: nil},
					auditortest.Expected{Action: createIdentityAction, Result: audit.Successful, Params: common.Params{"id": ID}},
				)
			})
		})
	})
}

func TestCreateIdentity_IdentityServiceError(t *testing.T) {
	Convey("should return expected error if identityService returns an error", t, func() {
		serviceMock := &apitest.IdentityServiceMock{
			CreateFunc: func(ctx context.Context, i *identity.Model) (string, error) {
				return "", errTest
			},
		}

		identityAPI := &API{IdentityService: serviceMock}

		newIdentity := &identity.Model{Name: "Eleven"}
		b, err := json.Marshal(newIdentity)
		So(err, ShouldBeNil)

		r := httptest.NewRequest("POST", createIdentityURL, bytes.NewReader(b))
		i, err := identityAPI.createIdentity(context.Background(), r)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, errTest)
		So(i, ShouldBeNil)
	})
}

func TestCreateIdentity_ReadBodyErr(t *testing.T) {
	Convey("should return expected error if reading the request body returns an error", t, func() {
		serviceMock := &apitest.IdentityServiceMock{}
		identityAPI := &API{IdentityService: serviceMock}

		r := httptest.NewRequest("POST", createIdentityURL, &IOReaderErroring{err: errors.New("")})
		i, err := identityAPI.createIdentity(context.Background(), r)
		So(err, ShouldEqual, ErrFailedToReadRequestBody)
		So(i, ShouldBeNil)
		So(serviceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}

func TestCreateIdentity_BodyEmpty(t *testing.T) {
	Convey("should return expected error if the request body is empty", t, func() {
		serviceMock := &apitest.IdentityServiceMock{}
		identityAPI := &API{IdentityService: serviceMock}

		r := httptest.NewRequest("POST", createIdentityURL, bytes.NewReader([]byte{}))
		i, err := identityAPI.createIdentity(context.Background(), r)
		So(err, ShouldEqual, ErrRequestBodyNil)
		So(i, ShouldBeNil)
		So(serviceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}

func TestAPI_CreateIdentityHandlerEmailAlreadyInUse(t *testing.T) {
	Convey("should return bad request if email already in use", t, func() {
		serviceMock := &apitest.IdentityServiceMock{
			CreateFunc: func(ctx context.Context, i *identity.Model) (string, error) {
				return "", identity.ErrEmailAlreadyExists
			},
		}

		auitorMock := auditortest.New()

		identityAPI := &API{
			IdentityService: serviceMock,
			auditor:         auitorMock,
		}

		newIdentity := &identity.Model{
			Name:  "Jamie",
			Email: "JamieLannister@GOT.com",
		}
		b, err := json.Marshal(newIdentity)
		So(err, ShouldBeNil)

		r := httptest.NewRequest("POST", createIdentityURL, bytes.NewReader(b))
		w := httptest.NewRecorder()

		identityAPI.CreateIdentityHandler(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
		So(strings.TrimSpace(w.Body.String()), ShouldEqual, identity.ErrEmailAlreadyExists.Error())

		So(serviceMock.CreateCalls(), ShouldHaveLength, 1)
		So(serviceMock.CreateCalls()[0].I.Name, ShouldEqual, "Jamie")
		So(serviceMock.CreateCalls()[0].I.Email, ShouldEqual, "JamieLannister@GOT.com")

		auitorMock.AssertRecordCalls(
			auditortest.Expected{Action: createIdentityAction, Result: audit.Attempted, Params: nil},
			auditortest.Expected{Action: createIdentityAction, Result: audit.Unsuccessful, Params: nil},
		)
	})
}

func assertErrorResponse(actualStatus int, expectedStatus int, actualBody string, expectedBody string) {
	So(actualStatus, ShouldEqual, expectedStatus)
	So(strings.TrimSpace(actualBody), ShouldEqual, expectedBody)
}
