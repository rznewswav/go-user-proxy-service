package t

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ms"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
)

var localeEn = en.New()
var localeZh = zh.New()
var localeMs = ms.New()
var uni = ut.New(localeEn, localeEn, localeZh, localeMs)

func GetTranslator(locale string) (trans ut.Translator) {
	trans, _ = uni.GetTranslator(locale)
	return
}

func init() {
	if enTranslator, found := uni.GetTranslator("en"); found {
		RegisterEnTranslations(enTranslator)
	}

	if msTranslator, found := uni.GetTranslator("ms"); found {
		RegisterMsTranslations(msTranslator)
	}

	if zhTranslator, found := uni.GetTranslator("zh"); found {
		RegisterZhTranslations(zhTranslator)
	}
}

//goland:noinspection GoUnhandledErrorResult
func RegisterEnTranslations(translator ut.Translator) {
	translator.Add("field-required", "\"{0}\" is required", true)
	translator.Add("number", "\"{0}\" must be a valid number", true)
}

//goland:noinspection GoUnhandledErrorResult
func RegisterMsTranslations(translator ut.Translator) {
	translator.Add("field-required", "\"{0}\" adalah diperlukan", true)
	translator.Add("number", "\"{0}\" mestilah jenis nombor", true)
}

//goland:noinspection GoUnhandledErrorResult,GoUnusedParameter
func RegisterZhTranslations(translator ut.Translator) {
	//translator.Add("field-required", "{0} is required", true)
}
