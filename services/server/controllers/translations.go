package controllers

import ut "github.com/go-playground/universal-translator"

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
