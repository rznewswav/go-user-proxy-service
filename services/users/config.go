package users

type NewswavUserConfig struct {
	NwApiBaseUrl string `env:"NW_API_BASE_URL" default:"https://dev-newswav-api.newswav.dev" required:"true"`
}
