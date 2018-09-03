package identity

import (
	"context"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	newIdentity = &Model{
		Name:  "Eleven",
		Email: "11@StrangerThings.com",
	}
)

func TestCreateIdentity_Success(t *testing.T) {
	Convey("should return no error if successful", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) error {
				return nil
			},
		}

		s := &Service{Persistence: persistenceMock}

		err := s.Create(context.Background(), newIdentity)

		So(err, ShouldBeNil)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 1)
		So(persistenceMock.CreateCalls()[0].Identity, ShouldResemble, newIdentity)
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

		err := s.Create(context.Background(), newIdentity)

		So(err, ShouldEqual, ErrPersistence)
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
