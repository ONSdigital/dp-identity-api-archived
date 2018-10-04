package identity

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//Create create a new user identity, returns ErrInvalidArguments if the required fields are invalid,
// ErrEmailAlreadyExists if the i.Email is already associated with an active identity and ErrPersistence for any errors
// persisting the identity to the datastore.
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

// VerifyPassword verify if the password provided is correct the email provided. Returns the active identity with the
// provided email if successful, returns ErrAuthenticateFailed if password incorrect.
func (s *Service) VerifyPassword(ctx context.Context, email string, password string) (*schema.Identity, error) {
	i, err := s.getIdentity(ctx, email)
	if err != nil {
		return nil, err
	}

	logD := log.Data{"email": email, "identity_id": i.ID}

	err = s.Encryptor.CompareHashAndPassword([]byte(i.Password), []byte(password))
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "password did not match stored value"), logD)
		return nil, ErrAuthenticateFailed
	}

	log.InfoCtx(ctx, "user authentication successful", logD)
	return i, nil
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
