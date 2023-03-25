package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh_tw"
	"reflect"
	"strings"
)

var (
	utrans *ut.UniversalTranslator
)

func InitTranslator() {
	v, _ := binding.Validator.Engine().(*validator.Validate)
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("comment"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	enTranslator := en.New()
	zhTwTranslator := zh_Hant_TW.New()
	utrans = ut.New(zhTwTranslator, zhTwTranslator, enTranslator)
	enTrans, _ := utrans.GetTranslator("en")
	zhTwTrans, _ := utrans.GetTranslator("zh_tw")
	en_translations.RegisterDefaultTranslations(v, enTrans)
	zh_translations.RegisterDefaultTranslations(v, zhTwTrans)
}

func GetValidationErrors(locale string, model interface{}, errors validator.ValidationErrors) map[string]string {
	validationErrors := make(map[string]string, len(errors))
	structType := reflect.TypeOf(model).Elem()
	for _, err := range errors {
		fieldName := err.StructField()
		field, _ := structType.FieldByName(fieldName)

		var fieldNameLoc string
		switch locale {
		case "en":
			fieldNameLoc = field.Tag.Get("en_field")
		case "zh_tw":
			fieldNameLoc = field.Tag.Get("zh_field")
		default:
			fieldNameLoc = field.Tag.Get("json")
		}

		tag := field.Tag.Get("json")
		validationErrors[tag] = strings.Replace(err.Translate(getTransFromParam(locale)), fieldName, fieldNameLoc, 1)
	}

	return validationErrors
}

func getTransFromParam(locale string) ut.Translator {
	translator, found := utrans.GetTranslator(locale)
	if !found {
		translator, _ = utrans.GetTranslator("zh_tw")
	}
	return translator
}
