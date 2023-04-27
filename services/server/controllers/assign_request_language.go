package controllers

const RequestLanguage = "mainlanguage"

var AssignRequestLanguage Handler[any] = func(
	Request Request[any],
	SetStatus SetStatus,
	SetHeader SetHeader,
) (Response any) {
	mainLangInHeader := Request.Header(RequestLanguage)
	if len(mainLangInHeader) == 0 {
		mainLangInHeader = "en"
	}

	Request.Set(RequestLanguage, mainLangInHeader)
	return
}
