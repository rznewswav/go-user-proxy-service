package users

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"service/services/bugsnag"
	"service/services/config"
	"service/services/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var nwApiGetUserInfoEndpoint string

func init() {
	newswavUserConfig := config.QuietBuild(NewswavUserConfig{})
	nwApiBaseUrl, urlParseError := url.Parse(newswavUserConfig.NwApiBaseUrl)
	if urlParseError != nil {
		panic(urlParseError)
	}

	joinedUrl, urlJoinError := url.JoinPath(nwApiBaseUrl.String(), "/v4/api/v1/me")
	if urlJoinError != nil {
		panic(urlJoinError)
	}

	nwApiGetUserInfoEndpoint = joinedUrl
}

var GetUserProfile = func(nwToken string) (
	success bool,
	profile AppUserType,
) {
	logger := logger.For("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request, requestInitError := http.NewRequestWithContext(
		ctx,
		"GET",
		nwApiGetUserInfoEndpoint,
		nil,
	)
	if requestInitError != nil {
		logger.Error(
			"error at sending integration request to %s",
			nwApiGetUserInfoEndpoint,
			bugsnag.FromError("User Profile Request Error", requestInitError),
		)
		success = false
		return
	}

	request.Header.Add("platform", "internal")
	request.Header.Add("nwtoken", nwToken)

	response, requestError := http.DefaultClient.Do(request)
	if requestError != nil {
		logger.Error(
			"error at sending integration request to %s",
			nwApiGetUserInfoEndpoint,
			bugsnag.FromError("User Profile Request Error", requestError),
		)
		success = false
		return
	}

	if response.StatusCode != http.StatusOK {
		logger.Error(
			"error at sending integration request to %s: %s",
			nwApiGetUserInfoEndpoint,
			response.Body,
			bugsnag.FromError("User Profile Request Error", requestError),
		)
		success = false
		return
	}

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		logger.Error(
			"error at reading integration response to %s",
			nwApiGetUserInfoEndpoint,
			bugsnag.FromError("User Profile Response Body Error", requestError),
		)
		success = false
		return
	}

	type User struct {
		MainLanguage     string   `json:"mainLanguage"`
		SubLanguages     []string `json:"subLanguages"`
		PnFrequency      int      `json:"pnFrequency"`
		LoginDisplayName string   `json:"loginDisplayName"`
		UserId           string   `json:"userId"`
		ProfilePicUrl    string   `json:"profilePicUrl"`
	}

	var user User
	unmarshalError := json.Unmarshal(body, &user)
	if unmarshalError != nil {
		logger.Error(
			"error at parsing integration response to json from %s",
			nwApiGetUserInfoEndpoint,
			bugsnag.FromError("User Profile Response Body Error", unmarshalError),
		)
		success = false
		return
	}

	type NewswavToken struct {
		jwt.RegisteredClaims
		Iss  string `json:"iss"`
		Aud  string `json:"aud"`
		Iat  int64  `json:"iat"`
		Exp  int64  `json:"exp"`
		PrId string `json:"pr_id"`
		Sdk  string `json:"sdk"`
		NwId int64  `json:"nw_id"`
		FiId string `json:"fi_id"`
		UId  string `json:"u_id"`
	}

	parsedJwt, jwtParsingError := jwt.ParseWithClaims(nwToken, &NewswavToken{}, nil)
	if jwtParsingError != nil && !errors.Is(jwtParsingError, jwt.ErrTokenUnverifiable) {
		logger.Error(
			"error at parsing JWT: %s",
			nwToken,
			bugsnag.FromError("User Profile JWT Error", requestError),
		)
		success = false
		return
	}

	if claims, ok := parsedJwt.Claims.(*NewswavToken); !ok {
		logger.Error(
			"error at parsing JWT: %s, JWT payload structure is not valid",
			nwToken,
		)
		success = false
		return
	} else {
		profile = AppUserType{
			MainLanguage:     user.MainLanguage,
			SubLanguages:     user.SubLanguages,
			PNFrequency:      user.PnFrequency,
			LoginDisplayName: user.LoginDisplayName,
			UserID:           claims.UId,
			ProfilePicURL:    user.ProfilePicUrl,
			ProfileID:        claims.PrId,
			NewswavID:        claims.NwId,
			FirebaseID:       claims.FiId,
		}

		success = true
		return

	}
}
