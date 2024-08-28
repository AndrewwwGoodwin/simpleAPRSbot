package api

import (
	"simpleAPRSbot-go/helpers/api/OpenWeatherMapWrapper"
	AprsFiWrapper "simpleAPRSbot-go/helpers/api/aprsFiWrapper"
	"simpleAPRSbot-go/helpers/api/osu"
)

type Keys struct {
	APRSFIkey         string
	OpenWeatherMapKey string
	OsuClientID       int
	OsuClientSecret   string
}

type Clients struct {
	APRSFi               *AprsFiWrapper.AprsFiClient
	OpenWeatherMapClient *OpenWeatherMapWrapper.OpenWeatherMapClient
	OSUClient            *osu.OsuAPIClient
}

func InitializeAPIClients(apiKeys Keys) Clients {
	var APRSFIClient = AprsFiWrapper.InitializeAprsFiClient(apiKeys.APRSFIkey)
	var OpenWeatherMapClient = OpenWeatherMapWrapper.New(apiKeys.OpenWeatherMapKey)
	var OsuClient, _ = osu.InitializeOsuClient(apiKeys.OsuClientID, apiKeys.OsuClientSecret, "client_credentials")
	return Clients{APRSFi: APRSFIClient,
		OpenWeatherMapClient: OpenWeatherMapClient,
		OSUClient:            OsuClient,
	}
}
