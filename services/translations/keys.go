package t

type TranslationKey string

const (
	GenericErrorTitle   TranslationKey = "GenericErrorTitle"
	GenericErrorMessage TranslationKey = "GenericErrorMessage"
	NotAuthorizedTitle  TranslationKey = "NotAuthorizedTitle"
	NotAuthorized       TranslationKey = "NotAuthorized"
)

func init() {
	AddTranslations(GenericErrorTitle, "Server encountered an error", "", "")
	AddTranslations(GenericErrorMessage, "Server encountered an error", "", "")
	AddTranslations(NotAuthorizedTitle, "Forbidden", "", "")
	AddTranslations(NotAuthorized, "You are not allowed to perform this operation", "", "")
}
