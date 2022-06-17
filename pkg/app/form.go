package app

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidError struct {
	Key     string
	Message string
}

func (v *ValidError) Error() string {
	return v.Message
}

type ValidErrors []*ValidError

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, error) {
	if err := c.ShouldBind(v); err != nil {
		return false, errors.Wrapf(err, "bind body fai")
	}

	validate, trans := initValidateAndTrans(c)
	if err := validate.Struct(v); err != nil {
		var errs ValidErrors
		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return false, errs
		}

		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}
		return false, errs
	}
	return true, nil
}

func initValidateAndTrans(c *gin.Context) (*validator.Validate, ut.Translator) {
	validate := validator.New()
	uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
	locale := c.GetHeader("locale") // 获取 local 参数
	trans, _ := uni.GetTranslator(locale)
	switch locale { // 列举 en 和 zh 的语言
	case "zh":
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	case "en":
		_ = en_translations.RegisterDefaultTranslations(validate, trans)
	default:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	}
	return validate, trans

}
