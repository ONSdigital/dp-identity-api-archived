package api

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"context"
	"net/http/httptest"
	"github.com/ONSdigital/dp-identity-api/models"
	"encoding/json"
	"bytes"
	"github.com/ONSdigital/dp-identity-api/store/storetest"
	"github.com/ONSdigital/dp-identity-api/store"
	"github.com/pkg/errors"
)

func TestCreateIdentity_Success(t *testing.T) {
	Convey("should return no error if successful", t, func() {
		mockStore := &storetest.StorerMock{
			CreateIdentityFunc: func(identity *models.Identity) error {
				return nil
			},
		}

		identityAPI := &IdentityAPI{
			dataStore: store.DataStore{Backend: mockStore},
		}

		newIdentity := &models.Identity{Name: "Eleven"}
		b, _ := json.Marshal(newIdentity)

		r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))

		err := identityAPI.createIdentity(context.Background(), r)

		So(err, ShouldBeNil)
		So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 1)
		So(mockStore.CreateIdentityCalls()[0].Identity, ShouldResemble, newIdentity)
	})
}

func TestCreateIdentity_ErrorUnmarshallingTheRequestBody(t *testing.T) {
	Convey("should return expected error if unmarshalling the request body returns an error", t, func() {
		mockStore := &storetest.StorerMock{
			CreateIdentityFunc: func(identity *models.Identity) error {
				return nil
			},
		}

		identityAPI := &IdentityAPI{
			dataStore: store.DataStore{Backend: mockStore},
		}

		b, _ := json.Marshal("not an identity")
		r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))

		err := identityAPI.createIdentity(context.Background(), r)

		So(err, ShouldEqual, ErrFailedToMarshalRequestBody)
		So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 0)
	})
}

func TestCreateIdentity_DataStoreError(t *testing.T) {
	Convey("should return expected error if datatore.CreateIdentity returns an error", t, func() {
		mockStore := &storetest.StorerMock{
			CreateIdentityFunc: func(identity *models.Identity) error {
				return errors.New("expected")
			},
		}

		identityAPI := &IdentityAPI{
			dataStore: store.DataStore{Backend: mockStore},
		}

		newIdentity := &models.Identity{Name: "Eleven"}
		b, _ := json.Marshal(newIdentity)
		r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))

		err := identityAPI.createIdentity(context.Background(), r)

		So(err, ShouldEqual, ErrFailedToWriteToMongo)
		So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 1)
		So(mockStore.CreateIdentityCalls()[0].Identity, ShouldResemble, newIdentity)
	})
}

func TestCreateIdentity_MissingParameters(t *testing.T) {
	Convey("should return expected error if context parameter is nil", t, func() {
		mockStore := &storetest.StorerMock{
			CreateIdentityFunc: func(identity *models.Identity) error {
				return errors.New("expected")
			},
		}

		identityAPI := &IdentityAPI{dataStore: store.DataStore{Backend: mockStore}}
		err := identityAPI.createIdentity(nil, nil)

		So(err, ShouldEqual, ErrInvalidArguments)
		So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 0)
	})

	Convey("should return expected error if context parameter is nil", t, func() {
		mockStore := &storetest.StorerMock{
			CreateIdentityFunc: func(identity *models.Identity) error {
				return errors.New("expected")
			},
		}

		identityAPI := &IdentityAPI{dataStore: store.DataStore{Backend: mockStore}}
		err := identityAPI.createIdentity(context.Background(), nil)

		So(err, ShouldEqual, ErrInvalidArguments)
		So(mockStore.CreateIdentityCalls(), ShouldHaveLength, 0)
	})
}
