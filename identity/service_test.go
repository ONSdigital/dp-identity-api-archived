package identity

import (
	"context"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	newIdentity = &Model{
		Name:     "Eleven",
		Email:    "11@StrangerThings.com",
		Password: "WAFFLES",
	}

	ID = "666"
)

func TestCreateIdentity_Success(t *testing.T) {
	Convey("should return no error if successful", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) (string, error) {
				return ID, nil
			},
		}

		s := &Service{Persistence: persistenceMock}

		id, err := s.Create(context.Background(), newIdentity)

		So(err, ShouldBeNil)
		So(id, ShouldEqual, ID)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 1)
		So(persistenceMock.CreateCalls()[0].Identity, ShouldResemble, newIdentity)
	})
}

func TestCreateIdentity_DataStoreError(t *testing.T) {
	Convey("should return expected error if datatore.CreateIdentity returns an error", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) (string, error) {
				return "", errors.New("expected")
			},
		}

		s := &Service{Persistence: persistenceMock}
		id, err := s.Create(context.Background(), newIdentity)

		So(err, ShouldEqual, ErrPersistence)
		So(id, ShouldBeEmpty)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 1)
		So(persistenceMock.CreateCalls()[0].Identity, ShouldResemble, newIdentity)
	})
}

func TestCreateIdentity_ValidationError(t *testing.T) {
	Convey("should return expected error if validate returns an error", t, func() {
		persistenceMock := &PersistenceMock{}
		s := &Service{Persistence: persistenceMock}

		id, err := s.Create(context.Background(), &Model{})

		So(err, ShouldResemble, ErrNameValidation)
		So(id, ShouldBeEmpty)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}

func TestCreateIdentity_MissingParameters(t *testing.T) {
	Convey("should return expected error if context parameter is nil", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) (string, error) {
				return "", errors.New("expected")
			},
		}

		s := &Service{Persistence: persistenceMock}

		id, err := s.Create(nil, nil)

		So(err, ShouldEqual, ErrInvalidArguments)
		So(id, ShouldBeEmpty)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})

	Convey("should return expected error if context parameter is nil", t, func() {
		persistenceMock := &PersistenceMock{
			CreateFunc: func(identity *Model) (string, error) {
				return "", errors.New("expected")
			},
		}

		s := &Service{Persistence: persistenceMock}

		id, err := s.Create(context.Background(), nil)

		So(err, ShouldEqual, ErrInvalidArguments)
		So(id, ShouldBeEmpty)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}

func TestService_Validate(t *testing.T) {
	s := Service{}

	Convey("should not return error if identity is valid", t, func() {
		i := &Model{
			Name:     "Bucky O'Hare",
			Email:    "captain@TheRighteousIndignation.com",
			Password: "S.P.A.C.E",
		}

		err := s.Validate(i)
		So(err, ShouldBeNil)
	})

	Convey("should error if identity is nil", t, func() {
		err := s.Validate(nil)
		So(err, ShouldResemble, ErrIdentityNil)
	})

	Convey("should error if identity.name is nil", t, func() {
		err := s.Validate(&Model{})
		So(err, ShouldResemble, ErrNameValidation)
	})

	Convey("should error if identity.email is nil", t, func() {
		err := s.Validate(&Model{Name: "Bucky O'Hare"})
		So(err, ShouldResemble, ErrEmailValidation)
	})

	Convey("should error if identity.password is nil", t, func() {
		err := s.Validate(&Model{Name: "Bucky O'Hare", Email: "captain@TheRighteousIndignation.com"})
		So(err, ShouldResemble, ErrPasswordValidation)
	})
}
