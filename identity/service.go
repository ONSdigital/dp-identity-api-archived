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
	ErrEmailAlreadyExists = errors.New("active identity already exists with email")
	ErrIdentityNotFound   = errors.New("authentication unsuccessful identity not found")
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

func (s *Service) VerifyPassword(ctx context.Context, email string, password string) error {
	i, err := s.getIdentity(ctx, email)
	if err != nil {
		return err
	}

	logD := log.Data{"email": email}

	err = s.Encryptor.CompareHashAndPassword([]byte(i.Password), []byte(password))
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "password did not match stored value"), logD)
		return ErrAuthenticateFailed
	}

	log.InfoCtx(ctx, "user authentication successful", logD)
	return nil
}

func (s *Service) getIdentity(ctx context.Context, email string) (*mongo.Identity, error) {
	logD := log.Data{"email": email}

	i, err := s.Persistence.GetIdentity(email)
	if err != nil {
		if err == mongo.ErrNotFound {
			log.ErrorCtx(ctx, errors.New("user not found"), logD)
			return nil, ErrIdentityNotFound
		}

		log.ErrorCtx(ctx, errors.Wrap(err, "error getting identity from database"), logD)
		return nil, err
	}
	return &i, nil
}

func (s *Service) encryptPassword(i *Model) (string, error) {
	pwd, err := s.Encryptor.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), err
}
