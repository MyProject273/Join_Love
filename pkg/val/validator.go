package val

import (
	"errors"

	"github.com/MyProject273/Join_Love/pkg/i18n"
	"github.com/go-playground/validator/v10"
)

func ParseValidationErrorI18n(err error, lang string) map[string]string {
	var verrs validator.ValidationErrors
	result := make(map[string]string)

	if errors.As(err, &verrs) {
		for _, fe := range verrs {
			field := fe.Field()
			msgID := "validation_" + fe.Tag()

			result[field] = i18n.GetI18nMessageWithData(lang, msgID, map[string]interface{}{
				"Field": field,
				"Param": fe.Param(),
			})
		}
	} else {
		result["general"] = i18n.GetI18nMessage(lang, "validation_invalid")
	}

	return result
}
