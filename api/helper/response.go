package helper

import (
	"net/http"

	"github.com/MyProject273/Join_Love/pkg/i18n"
	"github.com/MyProject273/Join_Love/pkg/utils/rescode"
	"github.com/MyProject273/Join_Love/pkg/val"
	"github.com/gin-gonic/gin"
)

// Struct for response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

// Function response return success
func RespondSuccess(ctx *gin.Context, data interface{}, code rescode.Code) {
	ctx.JSON(rescode.Success.HTTPStatus(), Response{
		Code:    code.Code(),
		Data:    data,
		Message: i18n.GetI18nMessage(rescode.Success.MessageKey(), GetLang(ctx)),
		Error:   nil,
	})
}

// Function response return error
func RespondError(ctx *gin.Context, code rescode.Code, err error) {
	ctx.JSON(code.HTTPStatus(), Response{
		Code:    code.Code(),
		Data:    nil,
		Message: i18n.GetI18nMessage(code.MessageKey(), GetLang(ctx)),
		Error:   err.Error(),
	})
}

func AbortWithErrorResponse(ctx *gin.Context, code rescode.Code, err error) {
	ctx.AbortWithStatusJSON(code.HTTPStatus(), Response{
		Code:    code.Code(),
		Data:    nil,
		Message: i18n.GetI18nMessage(code.MessageKey(), GetLang(ctx)),
		Error:   err.Error(),
	})
}

func RespondValidationError(ctx *gin.Context, err error) {
	errs := val.ParseValidationErrorI18n(err, GetLang(ctx))
	ctx.JSON(http.StatusBadRequest, Response{
		Code:    rescode.Invalid.Code(),
		Message: i18n.GetI18nMessage(rescode.Invalid.MessageKey(), GetLang(ctx)),
		Data:    nil,
		Error:   errs,
	})
}

func GetLang(ctx *gin.Context) string {
	lang := ctx.GetHeader("Accept-Language")
	if lang == "" {
		lang = "vi"
	}
	return lang
}
