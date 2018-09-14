package identity

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/mongo"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIdentityNil        = ValidationErr{message: "identity required but was nil"}
	ErrNameValidation     = ValidationErr{message: "mandatory field name was empty"}
	ErrEmailValidation    = ValidationErr{message: "mandatory field email was empty"}
	ErrPasswordValidation = ValidationErr{message: "mandatory field password was empty"}
	ErrAuthenticateFailed = errors.New("authentication unsuccessful")
	ErrUserNotFound       = errors.New("authentication unsuccessful user not found")
	ErrEmailAlreadyExists = errors.New("active identity already exists with email")
)

func (s *Service) Validate(i *Model) (err error) {
	if i == nil {
		return ErrIdentityNil
	}
	if i.Name == "" {
		return ErrNameValidation
	}
	if i.Email == "" {
		return ErrEmailValidation
	}
	if i.Password == "" {
		return ErrPasswordValidation
	}
	return nil
}

//Create create a new user identity
func (s *Service) Create(ctx context.Context, i *Model) (string, error) {
	if ctx == nil {
		log.Error(errors.New("create: failed mandatory context parameter was nil"), nil)
		return "", ErrInvalidArguments
	}

	if err := s.Validate(i); err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "create: failed validation"), nil)
		return "", err
	}

	logD := log.Data{
		"name":  i.Name,
		"email": i.Email,
	}

	pwd, err := s.encryptPassword(i)
	if err != nil {
		return "", errors.Wrap(err, "create: error encrypting password")
	}

	newIdentity := mongo.Identity{
		Name:              i.Name,
		Email:             i.Email,
		Password:          pwd,
		TemporaryPassword: i.TemporaryPassword,
		UserType:          i.UserType,
		Migrated:          i.Migrated,
		Deleted:           i.Deleted,
	}

	id, err := s.Persistence.Create(newIdentity)
	if err != nil && err == mongo.ErrNonUnique {
		log.ErrorCtx(ctx, errors.New("create: failed to create identity - an active identity with this email already exists"), logD)
		return "", ErrEmailAlreadyExists
	}

	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "create: failed to write data to mongo"), logD)
		return "", ErrPersistence
	}

	logD["id"] = id
	log.InfoCtx(ctx, "create: new identity created successfully", logD)
	return id, nil
}

func (s *Service) CreateToken(ctx context.Context, email string, password string) error {
	return nil
}

func (s *Service) encryptPassword(i *Model) (string, error) {
	pwd, err := s.Encryptor.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), err
}
