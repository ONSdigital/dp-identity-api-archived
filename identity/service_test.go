package identity

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/identity/identitytest"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/persistence/persistencetest"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

var (
	newIdentity = &schema.Identity{
		Name:     "Eleven",
		Email:    "11@StrangerThings.com",
		Password: "WAFFLES",
	}

	EMAIL = "666@ironmaiden.com"

	errTest = errors.New("test error")
)

func newPersistenceMock(email string, err error) *persistencetest.IdentityStoreMock {
	return &persistencetest.IdentityStoreMock{
		SaveIdentityFunc: func(identity schema.Identity) (string, error) {
			return email, err
		},
	}
}

func newEncryptorMock(pwd []byte, generate, compare error) *identitytest.EncryptorMock {
	return &identitytest.EncryptorMock{
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

		s := &Service{IdentityStore: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(context.Background(), newIdentity)

		So(err, ShouldBeNil)
		So(id, ShouldEqual, EMAIL)

		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(encryptorMock.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(encryptorMock.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)

		So(persistenceMock.SaveIdentityCalls(), ShouldHaveLength, 1)
		So(persistenceMock.SaveIdentityCalls()[0].NewIdentity, ShouldResemble, *newIdentity)
	})
}

func TestCreate_DataStoreError(t *testing.T) {
	Convey("should return expected error if datatore.CreateIdentity returns an error", t, func() {
		persistenceMock := newPersistenceMock("", errors.New("expected"))
		encryptorMock := newEncryptorMock([]byte(newIdentity.Password), nil, nil)

		s := &Service{IdentityStore: persistenceMock, Encryptor: encryptorMock}
		id, err := s.Create(context.Background(), newIdentity)

		So(err, ShouldEqual, ErrPersistence)
		So(id, ShouldBeEmpty)

		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(encryptorMock.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(encryptorMock.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)

		So(persistenceMock.SaveIdentityCalls(), ShouldHaveLength, 1)
		So(persistenceMock.SaveIdentityCalls()[0].NewIdentity, ShouldResemble, *newIdentity)
	})
}

func TestCreate_ValidationError(t *testing.T) {
	Convey("should return expected error if validate returns an error", t, func() {
		persistenceMock := &persistencetest.IdentityStoreMock{}
		encryptorMock := &identitytest.EncryptorMock{}
		s := &Service{IdentityStore: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(context.Background(), &schema.Identity{})

		So(err, ShouldResemble, schema.ErrNameValidation)
		So(id, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 0)
		So(persistenceMock.SaveIdentityCalls(), ShouldHaveLength, 0)
	})
}

func TestService_CreateEncryptPasswordError(t *testing.T) {
	Convey("should return expected error if encrypt password returns an error", t, func() {
		expectedErr := errors.New("encryption fails")

		persistenceMock := &persistencetest.IdentityStoreMock{}
		encryptorMock := newEncryptorMock([]byte{}, expectedErr, nil)

		s := &Service{IdentityStore: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(context.Background(), newIdentity)

		So(strings.Contains(err.Error(), "create: error encrypting password"), ShouldBeTrue)
		So(id, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(encryptorMock.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(encryptorMock.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)
		So(persistenceMock.SaveIdentityCalls(), ShouldHaveLength, 0)
	})
}

func TestCreate_MissingParameters(t *testing.T) {
	Convey("should return expected error if context parameter is nil", t, func() {
		persistenceMock := newPersistenceMock("", errors.New("expected"))
		encryptorMock := &identitytest.EncryptorMock{}

		s := &Service{IdentityStore: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(nil, nil)

		So(err, ShouldEqual, ErrInvalidArguments)
		So(id, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 0)
		So(persistenceMock.SaveIdentityCalls(), ShouldHaveLength, 0)
	})

	Convey("should return expected error if identity parameter is nil", t, func() {
		persistenceMock := newPersistenceMock("", errors.New("expected"))
		encryptorMock := &identitytest.EncryptorMock{}

		s := &Service{IdentityStore: persistenceMock, Encryptor: encryptorMock}

		id, err := s.Create(context.Background(), nil)

		So(err, ShouldResemble, schema.ErrIdentityNil)
		So(id, ShouldBeEmpty)
		So(encryptorMock.GenerateFromPasswordCalls(), ShouldHaveLength, 0)
		So(persistenceMock.SaveIdentityCalls(), ShouldHaveLength, 0)
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

func TestService_CreateEmailAlreadyInUse(t *testing.T) {
	Convey("todo", t, func() {
		persistenceMock := &persistencetest.IdentityStoreMock{
			SaveIdentityFunc: func(newIdentity schema.Identity) (string, error) {
				return "", persistence.ErrNonUnique
			},
		}

		enc := &identitytest.EncryptorMock{
			GenerateFromPasswordFunc: func(password []byte, cost int) ([]byte, error) {
				return []byte("WAFFLES"), nil
			},
		}

		s := Service{IdentityStore: persistenceMock, Encryptor: enc}

		_, err := s.Create(context.Background(), newIdentity)
		So(err, ShouldEqual, ErrEmailAlreadyExists)
		So(persistenceMock.SaveIdentityCalls(), ShouldHaveLength, 1)
		So(persistenceMock.SaveIdentityCalls()[0].NewIdentity, ShouldResemble, *newIdentity)
		So(enc.GenerateFromPasswordCalls(), ShouldHaveLength, 1)
		So(enc.GenerateFromPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(enc.GenerateFromPasswordCalls()[0].Cost, ShouldEqual, bcrypt.DefaultCost)

	})
}

func TestService_VerifyPassword(t *testing.T) {
	Convey("should not return error is password is correct", t, func() {
		p := &persistencetest.IdentityStoreMock{
			GetIdentityFunc: func(email string) (schema.Identity, error) {
				return *newIdentity, nil
			},
		}

		e := newEncryptorMock([]byte(newIdentity.Password), nil, nil)

		s := Service{IdentityStore: p, Encryptor: e}

		i, err := s.VerifyPassword(context.Background(), newIdentity.Email, newIdentity.Password)

		So(err, ShouldBeNil)
		So(i, ShouldResemble, newIdentity)
		So(p.GetIdentityCalls(), ShouldHaveLength, 1)
		So(p.GetIdentityCalls()[0].Email, ShouldEqual, newIdentity.Email)
		So(e.CompareHashAndPasswordCalls(), ShouldHaveLength, 1)
		So(e.CompareHashAndPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
		So(e.CompareHashAndPasswordCalls()[0].HashedPassword, ShouldResemble, []byte(newIdentity.Password))
	})
}

func TestService_VerifyPasswordIdentityNotFound(t *testing.T) {
	Convey("should return error if identity is not found", t, func() {
		p := &persistencetest.IdentityStoreMock{
			GetIdentityFunc: func(email string) (schema.Identity, error) {
				return schema.Identity{}, persistence.ErrNotFound
			},
		}

		e := newEncryptorMock([]byte(newIdentity.Password), nil, nil)

		s := Service{IdentityStore: p, Encryptor: e}

		i, err := s.VerifyPassword(context.Background(), newIdentity.Email, newIdentity.Password)

		So(err, ShouldEqual, ErrIdentityNotFound)
		So(i, ShouldBeNil)
		So(p.GetIdentityCalls(), ShouldHaveLength, 1)
		So(p.GetIdentityCalls()[0].Email, ShouldEqual, newIdentity.Email)
		So(e.CompareHashAndPasswordCalls(), ShouldHaveLength, 0)
	})
}

func TestService_VerifyPasswordPersistenceErr(t *testing.T) {
	Convey("should return error if identity is not found", t, func() {
		p := &persistencetest.IdentityStoreMock{
			GetIdentityFunc: func(email string) (schema.Identity, error) {
				return schema.Identity{}, errTest
			},
		}

		e := newEncryptorMock([]byte(newIdentity.Password), nil, nil)

		s := Service{IdentityStore: p, Encryptor: e}

		i, err := s.VerifyPassword(context.Background(), newIdentity.Email, newIdentity.Password)

		So(err, ShouldEqual, errTest)
		So(i, ShouldBeNil)
		So(p.GetIdentityCalls(), ShouldHaveLength, 1)
		So(p.GetIdentityCalls()[0].Email, ShouldEqual, newIdentity.Email)
		So(e.CompareHashAndPasswordCalls(), ShouldHaveLength, 0)
	})
}

func TestService_VerifyPasswordPasswordIncorrect(t *testing.T) {
	Convey("should return error if provided password is incorrect", t, func() {
		p := &persistencetest.IdentityStoreMock{
			GetIdentityFunc: func(email string) (schema.Identity, error) {
				return *newIdentity, nil
			},
		}

		e := newEncryptorMock([]byte(newIdentity.Password), nil, errTest)

		s := Service{IdentityStore: p, Encryptor: e}

		i, err := s.VerifyPassword(context.Background(), newIdentity.Email, newIdentity.Password)

		So(err, ShouldEqual, ErrAuthenticateFailed)
		So(i, ShouldBeNil)
		So(p.GetIdentityCalls(), ShouldHaveLength, 1)
		So(p.GetIdentityCalls()[0].Email, ShouldEqual, newIdentity.Email)
		So(e.CompareHashAndPasswordCalls(), ShouldHaveLength, 1)
		So(e.CompareHashAndPasswordCalls()[0].Password, ShouldResemble, []byte(newIdentity.Password))
	})
}
