package api

import (
	"simpleAPRSbot-go/helpers/api/OpenWeatherMap"
	"simpleAPRSbot-go/helpers/api/aprsFi"
	"simpleAPRSbot-go/helpers/api/osu"
)

type Keys struct {
	APRSFIkey         string
	OpenWeatherMapKey string
	OsuClientID       int
	OsuClientSecret   string
}

type Clients struct {
	APRSFi               *AprsFi.AprsFiClient
	OpenWeatherMapClient *OpenWeatherMap.OpenWeatherMapClient
	OSUClient            *osu.OsuAPIClient
}

func InitializeAPIClients(apiKeys Keys) Clients {
	var APRSFIClient = AprsFi.InitializeAprsFiClient(apiKeys.APRSFIkey)
	var OpenWeatherMapClient = OpenWeatherMap.New(apiKeys.OpenWeatherMapKey)
	var OsuClient, _ = osu.InitializeOsuClient(apiKeys.OsuClientID, apiKeys.OsuClientSecret, "client_credentials")
	return Clients{APRSFi: APRSFIClient,
		OpenWeatherMapClient: OpenWeatherMapClient,
		OSUClient:            OsuClient,
	}
}
