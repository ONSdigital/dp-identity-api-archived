package identity

import (
	"context"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
)

var (
	ErrIdentityNil        = ValidationErr{message: "identity required but was nil"}
	ErrNameValidation     = ValidationErr{message: "mandatory field name was empty"}
	ErrEmailValidation    = ValidationErr{message: "mandatory field email was empty"}
	ErrPasswordValidation = ValidationErr{message: "mandatory field password was empty"}
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
		log.Error(errors.New("Create: failed mandatory context parameter was nil"), nil)
		return "", ErrInvalidArguments
	}

	if i == nil {
		log.ErrorCtx(ctx, errors.New("Create: failed mandatory identity parameter was nil"), nil)
		return "", ErrInvalidArguments
	}

	s.Validate(i)

	id, err := s.Persistence.Create(i)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "Create: failed to write data to mongo"), nil)
		return "", ErrPersistence
	}

	log.InfoCtx(ctx, "Create: new identity created successfully", log.Data{
		"id":    id,
		"name":  i.Name,
		"email": i.Email,
	})
	return id, nil
}
