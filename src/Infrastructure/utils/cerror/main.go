package cerror

import "gitlab.com/yammt/oauth-auth-service/src/domain/auth/models"

func NewError(code int, err error) models.Error {
	return models.Error{
		Code: code,
		Err:  err,
	}
}
