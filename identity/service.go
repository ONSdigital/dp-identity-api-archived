package identity

import (
	"context"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
)

//Create create a new user identity
func (s *Service) Create(ctx context.Context, i *Model) error {
	if ctx == nil {
		log.Error(errors.New("Create: failed mandatory context parameter was nil"), nil)
		return ErrInvalidArguments
	}

	if i == nil {
		log.ErrorCtx(ctx, errors.New("Create: failed mandatory identity parameter was nil"), nil)
		return ErrInvalidArguments
	}

	err := s.Persistence.Create(i)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "Create: failed to write data to mongo"), nil)
		return ErrPersistence
	}

	log.InfoCtx(ctx, "Create: new identity created successfully", log.Data{"name": i.Name, "email": i.Email})
	return nil
}
