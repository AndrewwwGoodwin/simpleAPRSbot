package location

import (
	"github.com/ebarkie/aprs"
	"simpleAPRSbot-go/helpers/api"
	"simpleAPRSbot-go/helpers/api/OpenWeatherMapWrapper"
	"simpleAPRSbot-go/helpers/api/aprsFiWrapper"
	"simpleAPRSbot-go/helpers/aprsHelper"
	"strconv"
)

func Weather(args []string, f aprs.Frame, apiKey api.Keys) {
	if len(args) < 0 {
		// the user specified something! let's see what it was!
		aprsHelper.AprsTextReply("Not yet implemented!", f)
		return
	} else {
		// no location provided, lets just default to their last APRS location!
		var aprsfi = aprsFiWrapper.NewAprsFiWrapper(apiKey)
		var messageAuthor = aprsHelper.ExtractAuthor(f)

		locationInfo, err := aprsfi.GetLocation(messageAuthor)
		if err != nil {
			aprsHelper.AprsTextReply("Unable to get location", f)
			return
		}

		//with the location info, we need to give OpenWeatherMap a yell
		var owm = OpenWeatherMapWrapper.New(apiKey)
		err, weather := owm.GetWeather(locationInfo.Entries[0].Lat, locationInfo.Entries[0].Lng)
		if err != nil {
			aprsHelper.AprsTextReply("Unable to get weather", f)
			return
		}
		var messageToSend = weather.Daily[0].Summary + ". Currently it is " + weather.Current.Weather[0].Description + " " + strconv.FormatFloat(convertKtoF(weather.Current.Temp), 'f', 1, 64) + "F, Feels like " + strconv.FormatFloat(convertKtoF(weather.Current.FeelsLike), 'f', 1, 64) + "F, Humidity " + strconv.Itoa(weather.Current.Humidity) + "%"
		aprsHelper.AprsTextReply(messageToSend, f)
	}
}

func convertKtoF(input float64) float64 {
	return (input-273.15)*1.8 + 32
}
