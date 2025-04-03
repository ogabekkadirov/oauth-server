package cerror

import "github.com/ogabekkadirov/oauth-server/src/domain/auth/models"

func NewError(code int, err error) models.Error {
	return models.Error{
		Code: code,
		Err:  err,
	}
}
