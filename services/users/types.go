package users

type AppUserType struct {
	MainLanguage     string   `json:"mainLanguage"`
	SubLanguages     []string `json:"subLanguages"`
	PNFrequency      int      `json:"pnFrequency"`
	LoginDisplayName string   `json:"loginDisplayName"`
	UserID           string   `json:"userId"`
	ProfilePicURL    string   `json:"profilePicUrl"`
	ProfileID        string   `json:"profileId"`
	NewswavID        int64    `json:"newswavId"`
	FirebaseID       string   `json:"firebaseId"`
}
