package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ogabekkadirov/oauth-server/src/Infrastructure/helpers"
	"github.com/ogabekkadirov/oauth-server/src/domain/auth/models"
)

type MyStruct struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

type ResponseResult struct {
	Success      bool        `json:"success" xml:"success"`
	Message      string      `json:"message" xml:"message"`
	Result       interface{} `json:"result" xml:"result"`
	Status       int         `json:"status" xml:"status"`
	ErrorMessage interface{} `json:"error_message" xml:"error_message"`
}

type ValidationParams struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

type ResponseError struct {
	ErrorMessage     string              `json:"error_message"`
	Message          string              `json:"message"`
	Status           int16               `json:"status"`
	Success          bool                `json:"success" default:"false"`
	ValidationErrors []*ValidationParams `json:"validate_errors"`
}

type InternalServerError struct {
	ErrorMessage string `json:"error_message"`
	Message      string `json:"message" default:"Internal Server Error"`
	Status       int16  `json:"status" default:"500"`
	Success      bool   `json:"success" default:"false"`
}

func SuccessResult(ctx *gin.Context, httpStatus int, result interface{}) (mapResult map[string]interface{}) {

	mapResult = make(map[string]interface{})
	mapResult["success"] = true
	mapResult["message"] = "ok"
	mapResult["result"] = result
	mapResult["status"] = httpStatus

	ctx.JSON(httpStatus, mapResult)

	return
}

func ErrorResult(ctx *gin.Context, err models.Error) (mapResult map[string]interface{}) {

	mapResult = make(map[string]interface{})
	mapResult["success"] = false
	mapResult["message"] = http.StatusText(err.Code)
	mapResult["status"] = err.Code
	mapResult["error_message"] = err.Err.Error()
	if err.Code == 422 {
		mapResult["error_message"] = "validation error"
		var errData models.ErrBody
		_ = json.Unmarshal([]byte(err.Err.Error()), &errData)

		mapResult["validate_errors"] = errData.Error
	} else if err.Err != nil {
		var ve validator.ValidationErrors
		if errors.As(err.Err, &ve) {
			out := make([]models.ValidatorError, len(ve))
			for i, fe := range ve {
				out[i] = models.ValidatorError{Field: fe.Field(), Msg: helpers.MsgForTag(fe.Tag())}
			}
			println(out)
			mapResult["validate_errors"] = out
		}
	}

	ctx.AbortWithStatusJSON(err.Code, mapResult)

	return
}
