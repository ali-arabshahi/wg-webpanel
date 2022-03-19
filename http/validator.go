package http

import (
	"log"
	"net"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var modelValidator *validator.Validate

func newValidator() {
	modelValidator = validator.New()
	modelValidator.RegisterValidation("iplist", func(fl validator.FieldLevel) bool {
		ipList := fl.Field().Interface().([]string)
		if len(ipList) == 0 {
			return false
		}
		for _, ip := range ipList {
			_, _, err := net.ParseCIDR(ip)
			if err != nil {
				return false
			}
		}
		return true
	})
	//------------- add custome error -------------------//
	translator := en.New()
	uni := ut.New(translator, translator)
	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}
	//-----------------------------------------------------//
	if err := en_translations.RegisterDefaultTranslations(modelValidator, trans); err != nil {
		log.Fatal(err)
	}
	//---------------------------------------------------------
	_ = modelValidator.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	//---------------------------------------------------------
	_ = modelValidator.RegisterTranslation("ip", trans, func(ut ut.Translator) error {
		return ut.Add("ip", "{0} is not a valid ip E.g. x.x.x.x", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ip", fe.Field())
		return t
	})
	//---------------------------------------------------------
	_ = modelValidator.RegisterTranslation("cidr", trans, func(ut ut.Translator) error {
		return ut.Add("cidr", "{0} is not a valid ip cidr E.g. x.x.x.x/x", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("cidr", fe.Field())
		return t
	})
	//---------------------------------------------------------
	_ = modelValidator.RegisterTranslation("iplist", trans, func(ut ut.Translator) error {
		return ut.Add("iplist", "{0} is not a valid ip or list is empty", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("iplist", fe.Field())
		return t
	})

}
