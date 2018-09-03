package identity

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// TODO refactor Create to accept model instead of request.

//Create create a new user identity
func (s *Service) Create(ctx context.Context, r *http.Request) error {
	if ctx == nil {
		log.Error(errors.New("Create: failed mandatory context parameter was nil"), nil)
		return ErrInvalidArguments
	}

	if r == nil {
		log.Error(errors.New("Create: failed mandatory request parameter was nil"), nil)
		return ErrInvalidArguments
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "Create: failed to read request body"), nil)
		return ErrFailedToReadRequestBody
	}

	defer r.Body.Close()

	var i *Model
	err = json.Unmarshal(body, &i)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "Create: failed to unmarshal request body"), nil)
		return ErrFailedToUnmarshalRequestBody
	}

	err = s.Persistence.Create(i)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "Create: failed to write data to mongo"), nil)
		return ErrFailedToWriteToMongo
	}

	log.InfoCtx(ctx, "Create: new identity created successfully", log.Data{"name": i.Name, "email": i.Email})
	return nil
}
