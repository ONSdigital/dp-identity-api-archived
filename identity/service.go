package identity

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAuthenticateFailed = errors.New("authentication unsuccessful")
	ErrEmailAlreadyExists = errors.New("active identity already exists with email")
	ErrIdentityNotFound   = errors.New("authentication unsuccessful user not found")
)

//Create create a new user identity
func (s *Service) Create(ctx context.Context, i *schema.Identity) (string, error) {
	if ctx == nil {
		log.Error(errors.New("create: failed mandatory context parameter was nil"), nil)
		return "", ErrInvalidArguments
	}

	if err := i.Validate(); err != nil {
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

	i.Password = pwd

	id, err := s.IdentityStore.SaveIdentity(*i)
	if err != nil && err == persistence.ErrNonUnique {
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

func (s *Service) Get(ctx context.Context, tokenStr string) (*schema.Identity, error) {

	// TODO - has token expired?
	// TODO - token to get id from cache
	// TODO - id to get requested identity

	defaultUser := &schema.Identity{
		Name:              "John Paul Jones",
		Email:             "blackdog@ons.gov.uk",
		Password:          "foo",
		UserType:          "bar",
		TemporaryPassword: false,
		Migrated:          false,
		Deleted:           false,
	}

	return defaultUser, nil
}

func (s *Service) getIdentity(ctx context.Context, email string) (*schema.Identity, error) {
	logD := log.Data{"email": email}

	i, err := s.IdentityStore.GetIdentity(email)
	if err != nil {
		if err == persistence.ErrNotFound {
			log.ErrorCtx(ctx, errors.New("user not found"), logD)
			return nil, ErrIdentityNotFound
		}

		log.ErrorCtx(ctx, errors.Wrap(err, "error getting identity from database"), logD)
		return nil, err
	}
	return &i, nil

}

func (s *Service) encryptPassword(i *schema.Identity) (string, error) {
	pwd, err := s.Encryptor.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), err
}
