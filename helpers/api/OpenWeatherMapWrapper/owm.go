package OpenWeatherMapWrapper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"simpleAPRSbot-go/helpers/api"
)

type OpenWeatherMapWrapper struct {
	apikey string
}

func New(apikey api.Keys) *OpenWeatherMapWrapper {
	return &OpenWeatherMapWrapper{apikey: apikey.OpenWeatherMapkey}
}

func (owm OpenWeatherMapWrapper) GetWeather(lat string, lon string) (error, *WeatherJsonFormat) {
	//https://api.openweathermap.org/data/3.0/onecall?lat={lat}&lon={lon}&exclude={part}&appid={API key}
	resp, err := http.Get("https://api.openweathermap.org/data/3.0/onecall?lat=" + lat + "&lon=" + lon + "&exclude=minutely,hourly,alerts&appid=" + owm.apikey)
	if err != nil {
		return errors.New("error getting weather from OpenWeatherMap"), nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	var weatherJson WeatherJsonFormat
	err = json.Unmarshal(body, &weatherJson)
	if err != nil {
		return errors.New("error unmarshalling OpenWeatherMap"), nil
	}
	return nil, &weatherJson
}

type WeatherJsonFormat struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        struct {
		Dt         int     `json:"dt"`
		Sunrise    int     `json:"sunrise"`
		Sunset     int     `json:"sunset"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Uvi        float64 `json:"uvi"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		Weather    []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Daily []struct {
		Dt        int     `json:"dt"`
		Sunrise   int     `json:"sunrise"`
		Sunset    int     `json:"sunset"`
		Moonrise  int     `json:"moonrise"`
		Moonset   int     `json:"moonset"`
		MoonPhase float64 `json:"moon_phase"`
		Summary   string  `json:"summary"`
		Temp      struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float64 `json:"day"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		DewPoint  float64 `json:"dew_point"`
		WindSpeed float64 `json:"wind_speed"`
		WindDeg   int     `json:"wind_deg"`
		WindGust  float64 `json:"wind_gust"`
		Weather   []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds int     `json:"clouds"`
		Pop    float64 `json:"pop"`
		Rain   float64 `json:"rain,omitempty"`
		Uvi    float64 `json:"uvi"`
	} `json:"daily"`
}
