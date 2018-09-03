package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

func TestCreateIdentity_Success(t *testing.T) {
	Convey("should return no error if successful", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) error {
				return nil
			},
		}

		s := &Service{Persistence: persistenceMock}

		newIdentity := &Model{Name: "Eleven"}
		b, _ := json.Marshal(newIdentity)

		r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))

		err := s.Create(context.Background(), r)

		So(err, ShouldBeNil)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 1)
		So(persistenceMock.CreateCalls()[0].Identity, ShouldResemble, newIdentity)
	})
}

func TestCreateIdentity_ErrorUnmarshalingTheRequestBody(t *testing.T) {
	Convey("should return expected error if unmarshalling the request body returns an error", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) error {
				return nil
			},
		}

		s := &Service{Persistence: persistenceMock}

		b, _ := json.Marshal("not an identity")
		r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))

		err := s.Create(context.Background(), r)

		So(err, ShouldEqual, ErrFailedToUnmarshalRequestBody)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}

func TestCreateIdentity_DataStoreError(t *testing.T) {
	Convey("should return expected error if datatore.CreateIdentity returns an error", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) error {
				return errors.New("expected")
			},
		}

		s := &Service{Persistence: persistenceMock}

		newIdentity := &Model{Name: "Eleven"}
		b, _ := json.Marshal(newIdentity)
		r := httptest.NewRequest("POST", "http://localhost:23800/identity", bytes.NewReader(b))

		err := s.Create(context.Background(), r)

		So(err, ShouldEqual, ErrFailedToWriteToMongo)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 1)
		So(persistenceMock.CreateCalls()[0].Identity, ShouldResemble, newIdentity)
	})
}

func TestCreateIdentity_MissingParameters(t *testing.T) {
	Convey("should return expected error if context parameter is nil", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) error {
				return errors.New("expected")
			},
		}

		s := &Service{Persistence: persistenceMock}

		err := s.Create(nil, nil)

		So(err, ShouldEqual, ErrInvalidArguments)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})

	Convey("should return expected error if context parameter is nil", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) error {
				return errors.New("expected")
			},
		}

		s := &Service{Persistence: persistenceMock}

		err := s.Create(context.Background(), nil)

		So(err, ShouldEqual, ErrInvalidArguments)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}
