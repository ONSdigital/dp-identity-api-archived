package identity

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/mongo"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

var (
	newIdentity = &Model{
		Name:     "Eleven",
		Email:    "11@StrangerThings.com",
		Password: "WAFFLES",
	}

	newMongoIdentity = mongo.Identity{
		Name:     "Eleven",
		Email:    "11@StrangerThings.com",
		Password: "WAFFLES",
	}

	EMAIL = "666@ironmaiden.com"

	errTest = errors.New("mongo error")
)

func newPersistenceMock(email string, err error) *PersistenceMock {
	return &PersistenceMock{
		CreateFunc: func(identity mongo.Identity) (string, error) {
			return email, err
		},
	}
}

func newEncryptorMock(pwd []byte, generate, compare error) *EncryptorMock {
	return &EncryptorMock{
		GenerateFromPasswordFunc: func(password []byte, cost int) ([]byte, error) {
			return pwd, generate
		},

		CompareHashAndPasswordFunc: func(hashedPassword []byte, password []byte) error {
			return compare
		},
	}
}

func TestCreate_Success(t *testing.T) {
	Convey("should return no error if successful", t, func() {
		persistenceMock := newPersistenceMock(EMAIL, nil)
		encryptorMock := newEncryptorMock([]byte(newIdentity.Password), nil, nil)

		s := &Service{Persistence: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(context.Background(), newIdentity)

		So(err, ShouldBeNil)
		So(id, ShouldEqual, EMAIL)

		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(encryptorMock.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(encryptorMock.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)

		So(persistenceMock.CreateCalls(), ShouldHaveLength, 1)
		So(persistenceMock.CreateCalls()[0].NewIdentity, ShouldResemble, newMongoIdentity)
	})
}

func TestCreate_DataStoreError(t *testing.T) {
	Convey("should return expected error if datatore.CreateIdentity returns an error", t, func() {
		persistenceMock := newPersistenceMock("", errors.New("expected"))
		encryptorMock := newEncryptorMock([]byte(newIdentity.Password), nil, nil)

		s := &Service{Persistence: persistenceMock, Encryptor: encryptorMock}
		id, err := s.Create(context.Background(), newIdentity)

		So(err, ShouldEqual, ErrPersistence)
		So(id, ShouldBeEmpty)

		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(encryptorMock.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(encryptorMock.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)

		So(persistenceMock.CreateCalls(), ShouldHaveLength, 1)
		So(persistenceMock.CreateCalls()[0].NewIdentity, ShouldResemble, newMongoIdentity)
	})
}

func TestCreate_ValidationError(t *testing.T) {
	Convey("should return expected error if validate returns an error", t, func() {
		persistenceMock := &PersistenceMock{}
		encryptorMock := &EncryptorMock{}
		s := &Service{Persistence: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(context.Background(), &Model{})

		So(err, ShouldResemble, ErrNameValidation)
		So(id, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 0)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}

func TestService_CreateEncryptPasswordError(t *testing.T) {
	Convey("should return expected error if encrypt password returns an error", t, func() {
		expectedErr := errors.New("encryption fails")

		persistenceMock := &PersistenceMock{}
		encryptorMock := newEncryptorMock([]byte{}, expectedErr, nil)

		s := &Service{Persistence: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(context.Background(), newIdentity)

		So(strings.Contains(err.Error(), "create: error encrypting password"), ShouldBeTrue)
		So(id, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(encryptorMock.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(encryptorMock.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})
}

func TestCreate_MissingParameters(t *testing.T) {
	Convey("should return expected error if context parameter is nil", t, func() {
		persistenceMock := newPersistenceMock("", errors.New("expected"))
		encryptorMock := &EncryptorMock{}

		s := &Service{Persistence: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(nil, nil)

		So(err, ShouldEqual, ErrInvalidArguments)
		So(id, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 0)
		So(persistenceMock.CreateCalls(), ShouldHaveLength, 0)
	})

	Convey("should return expected error if identity parameter is nil", t, func() {
		persistenceMock := newPersistenceMock("", errors.New("expected"))
		encryptorMock := &EncryptorMock{}

		s := &Service{Persistence: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(context.Background(), nil)

		So(err, ShouldResemble, ErrIdentityNil)
		So(id, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 0)
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

func TestService_EncryptPassword(t *testing.T) {
	Convey("should return expected value if encryption is successful", t, func() {
		expectedPWD := "123456789"
		encryptorMock := newEncryptorMock([]byte(expectedPWD), nil, nil)

		s := Service{Encryptor: encryptorMock}

		pwd, err := s.encryptPassword(newIdentity)

		So(err, ShouldBeNil)
		So(pwd, ShouldResemble, expectedPWD)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(encryptorMock.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(encryptorMock.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)
	})

	Convey("should return expected error if encryption is unsuccessful", t, func() {
		expectedPWD := "123456789"
		expectedErr := errors.New("encryption")
		encryptorMock := newEncryptorMock([]byte(expectedPWD), expectedErr, nil)

		s := Service{Encryptor: encryptorMock}

		pwd, err := s.encryptPassword(newIdentity)

		So(err, ShouldEqual, expectedErr)
		So(pwd, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(encryptorMock.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(encryptorMock.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)
	})
}

func TestService_VerifyPassword(t *testing.T) {
	Convey("should not return error is password is correct", t, func() {
		p := &PersistenceMock{
			GetIdentityFunc: func(email string) (mongo.Identity, error) {
				return newMongoIdentity, nil
			},
		}

		e := newEncryptorMock([]byte(newMongoIdentity.Password), nil, nil)

		s := Service{Persistence: p, Encryptor: e}

		err := s.VerifyPassword(context.Background(), newMongoIdentity.Email, newMongoIdentity.Password)

		So(err, ShouldBeNil)
		So(p.GetIdentityCalls(), ShouldHaveLength, 1)
		So(p.GetIdentityCalls()[0].Email, ShouldEqual, newMongoIdentity.Email)
		So(e.CompareHashAndPasswordCalls(), ShouldHaveLength, 1)
		So(e.CompareHashAndPasswordCalls()[0].Password, ShouldResemble, []byte(newMongoIdentity.Password))
		So(e.CompareHashAndPasswordCalls()[0].HashedPassword, ShouldResemble, []byte(newMongoIdentity.Password))
	})
}

func TestService_VerifyPasswordIdentityNotFound(t *testing.T) {
	Convey("should return error if identity is not found", t, func() {
		p := &PersistenceMock{
			GetIdentityFunc: func(email string) (mongo.Identity, error) {
				return mongo.Identity{}, mongo.ErrNotFound
			},
		}

		e := newEncryptorMock([]byte(newMongoIdentity.Password), nil, nil)

		s := Service{Persistence: p, Encryptor: e}

		err := s.VerifyPassword(context.Background(), newMongoIdentity.Email, newMongoIdentity.Password)

		So(err, ShouldEqual, ErrIdentityNotFound)
		So(p.GetIdentityCalls(), ShouldHaveLength, 1)
		So(p.GetIdentityCalls()[0].Email, ShouldEqual, newMongoIdentity.Email)
		So(e.CompareHashAndPasswordCalls(), ShouldHaveLength, 0)
	})
}

func TestService_VerifyPasswordPersistenceErr(t *testing.T) {
	Convey("should return error if identity is not found", t, func() {
		p := &PersistenceMock{
			GetIdentityFunc: func(email string) (mongo.Identity, error) {
				return mongo.Identity{}, errTest
			},
		}

		e := newEncryptorMock([]byte(newMongoIdentity.Password), nil, nil)

		s := Service{Persistence: p, Encryptor: e}

		err := s.VerifyPassword(context.Background(), newMongoIdentity.Email, newMongoIdentity.Password)

		So(err, ShouldEqual, errTest)
		So(p.GetIdentityCalls(), ShouldHaveLength, 1)
		So(p.GetIdentityCalls()[0].Email, ShouldEqual, newMongoIdentity.Email)
		So(e.CompareHashAndPasswordCalls(), ShouldHaveLength, 0)
	})
}

func TestService_VerifyPasswordPasswordIncorrect(t *testing.T) {
	Convey("should return error if provided password is incorrect", t, func() {
		p := &PersistenceMock{
			GetIdentityFunc: func(email string) (mongo.Identity, error) {
				return newMongoIdentity, nil
			},
		}

		e := newEncryptorMock([]byte(newMongoIdentity.Password), nil, errTest)

		s := Service{Persistence: p, Encryptor: e}

		err := s.VerifyPassword(context.Background(), newMongoIdentity.Email, newMongoIdentity.Password)

		So(err, ShouldEqual, ErrAuthenticateFailed)
		So(p.GetIdentityCalls(), ShouldHaveLength, 1)
		So(p.GetIdentityCalls()[0].Email, ShouldEqual, newMongoIdentity.Email)
		So(e.CompareHashAndPasswordCalls(), ShouldHaveLength, 1)
		So(e.CompareHashAndPasswordCalls()[0].Password, ShouldResemble, []byte(newMongoIdentity.Password))
	})
}
