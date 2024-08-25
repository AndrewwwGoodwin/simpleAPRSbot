package osu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type AuthDataType struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type OsuAPIClient struct {
	grantType    string
	clientId     int
	clientSecret string
	ccToken      string
	expiration   time.Time
	httpClient   *http.Client
}

type ModeString string

var (
	modeOsuValue   ModeString = "osu"
	modeTaikoValue ModeString = "taiko"
	modeCTBValue   ModeString = "ctb"
	modeManiaValue ModeString = "mania"

	ModeOsu   = &modeOsuValue
	ModeTaiko = &modeTaikoValue
	ModeCTB   = &modeCTBValue
	ModeMania = &modeManiaValue
)

func InitializeOsuClient(clientId int, clientSecret string, authType string) (*OsuAPIClient, error) {
	// add check that authType is client_credentials or whatever the other one is and switch
	var httpClient = &http.Client{}
	// ensure that clientID and clientSecret make sense too
	if clientId == 0 || clientSecret == "" {
		return nil, errors.New("clientId and clientSecret are required")
	}
	return &OsuAPIClient{clientId: clientId, clientSecret: clientSecret, grantType: authType, httpClient: httpClient}, nil
}

func (client *OsuAPIClient) Authenticate() error {
	if time.Now().Before(client.expiration) {
		return nil
	}

	var authURL, _ = url.Parse("https://osu.ppy.sh/oauth/token")

	params := url.Values{}
	params.Add("client_id", strconv.Itoa(client.clientId))
	params.Add("client_secret", client.clientSecret)
	params.Add("grant_type", client.grantType)
	params.Add("scope", "public")

	authURL.RawQuery = params.Encode()

	req, err := http.PostForm(authURL.String(), params)
	if err != nil {
		return fmt.Errorf("failed to post osu oauth %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(req.Body)

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	var formattedData AuthDataType

	err = json.Unmarshal(data, &formattedData)
	if err != nil {
		return fmt.Errorf("failed to Unmarshal Auth Data: %s", err)
	}
	expireDuration := time.Duration(formattedData.ExpiresIn-30) * time.Second
	expirationTime := time.Now().Add(expireDuration)
	client.expiration = expirationTime
	client.ccToken = formattedData.AccessToken
	return nil
}
