package api

import (
	"testing"
	"github.com/ONSdigital/go-ns/audit/auditortest"
	"github.com/ONSdigital/go-ns/audit"
	"net/http/httptest"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"github.com/ONSdigital/dp-identity-api/models"
	"errors"
	"github.com/ONSdigital/dp-identity-api/store/storetest"
	"github.com/ONSdigital/dp-identity-api/store"
)

func TestIdentityAPI_CreateIdentityHandler(t *testing.T) {
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
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
				So(strings.TrimSpace(w.Body.String()), ShouldEqual, auditortest.ErrAudit.Error())
			})

			Convey("and no identity is created", func() {
				So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 0)
			})
		})
	})
}
