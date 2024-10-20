package api

import (
	"fmt"
	"simpleAPRSbot-go/helpers/api/OpenWeatherMap"
	"simpleAPRSbot-go/helpers/api/aprsFi"
	"simpleAPRSbot-go/helpers/api/osu"
)

type Keys struct {
	APRSFIkey         *string
	OpenWeatherMapKey *string
	OsuClientID       *int
	OsuClientSecret   *string
}

type Clients struct {
	APRSFi               *AprsFi.AprsFiClient
	OpenWeatherMapClient *OpenWeatherMap.OpenWeatherMapClient
	OSUClient            *osu.OsuAPIClient
}

func InitializeAPIClients(apiKeys *Keys) Clients {
	var returnObject = Clients{}
	//first we need to check if the provided keys are nil
	if apiKeys.OpenWeatherMapKey != nil {
		// we nest this so that we don't get a nil-pointer deref
		if *apiKeys.OpenWeatherMapKey != "" {
			returnObject.OpenWeatherMapClient = OpenWeatherMap.New(*apiKeys.OpenWeatherMapKey)
		}
	}
	if apiKeys.APRSFIkey != nil {
		if *apiKeys.APRSFIkey != "" {
			returnObject.APRSFi = AprsFi.InitializeAprsFiClient(*apiKeys.APRSFIkey)
		}
	}
	if apiKeys.OsuClientID != nil || apiKeys.OsuClientSecret != nil {
		var OSUClient *osu.OsuAPIClient
		var err error
		if *apiKeys.OsuClientID != 0 && *apiKeys.OsuClientSecret != "" {
			OSUClient, err = osu.InitializeOsuClient(*apiKeys.OsuClientID, *apiKeys.OsuClientSecret, "client_credentials")
			if err != nil {
				fmt.Println("Error initializing Osu Client")
				returnObject.OSUClient = nil
			}
		}

		returnObject.OSUClient = OSUClient
	}
	return returnObject
}
