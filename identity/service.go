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
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "create: failed to write data to mongo"), nil)
		return "", ErrPersistence
	}

	log.InfoCtx(ctx, "create: new identity created successfully", log.Data{
		"id":    id,
		"name":  i.Name,
		"email": i.Email,
	})
	return id, nil
}

func (s *Service) Authenticate(ctx context.Context, id string, password string) error {
	i, err := s.getIdentity(ctx, id)
	if err != nil {
		return err
	}

	logD := log.Data{"id": id}

	err = s.Encryptor.CompareHashAndPassword([]byte(i.Password), []byte(password))
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "password did not match stored value"), logD)
		return ErrAuthenticateFailed
	}

	log.InfoCtx(ctx, "user authentication successful", logD)
	return nil
}

func (s *Service) getIdentity(ctx context.Context, id string) (*mongo.Identity, error) {
	logD := log.Data{"id": id}

	i, err := s.Persistence.GetIdentity(id)
	if err != nil {
		if err == mongo.ErrNotFound {
			log.ErrorCtx(ctx, errors.New("user not found"), logD)
			return nil, ErrUserNotFound
		}

		log.ErrorCtx(ctx, errors.Wrap(err, "error getting identity from database"), logD)
		return nil, err
	}
	return i, nil
}

func (s *Service) encryptPassword(i *Model) (string, error) {
	pwd, err := s.Encryptor.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), err
}
