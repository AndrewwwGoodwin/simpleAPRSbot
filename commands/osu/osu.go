package osuCommands

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"math"
	"simpleAPRSbot-go/helpers/APRS"
	"simpleAPRSbot-go/helpers/api/osu"
	"strconv"
	"strings"
)

// copy BathBot behavior and return a basic summary of a user's profile

func Osu(args []string, f aprs.Frame, client *APRS.UserClient) {
	var username = args[0]
	if username == "" {
		client.Reply("You need to specify a username!", f)
	}

	var osuClient = client.ApiClients.OSUClient

	user, err := osuClient.GetUser(username, osu.ModeOsu, osu.KeyUsername)
	if err != nil {
		client.Reply("Error getting user", f)
		return
	}
	if user == nil {
		client.Reply("User not found", f)
		return
	}
	fmt.Println(user)
	var userUsername = user.Username
	var userPP = FloatToCommaString(user.Statistics.Pp) + "pp"
	var userGlobalRank = "#" + IntToCommaString(user.Statistics.GlobalRank)
	var userCountryRank = user.CountryCode + IntToCommaString(user.Statistics.CountryRank)
	var userAcc = "Accuracy:" + strconv.FormatFloat(user.Statistics.HitAccuracy, 'f', 2, 64) + "%"
	var userLevel = "Level: " + strconv.Itoa(user.Statistics.Level.Current) + "." + strconv.Itoa(user.Statistics.Level.Progress)
	var userPlaycount = "Playcount: " + IntToCommaString(user.Statistics.PlayCount)
	var userHours = IntToCommaString(int(math.Round(float64((user.Statistics.PlayTime/60)/60)))) + "hrs"
	var userMedalCount = "Medals: " + strconv.Itoa(len(user.UserAchievements))
	var userPeakRank = "Peak Rank: #" + IntToCommaString(user.RankHighest.Rank) + " (" + user.RankHighest.UpdatedAt.Format("01/02/2006") + ")"
	var returnString = userUsername + ": " + userPP + " (" + userGlobalRank + " " + userCountryRank + ") " + userAcc + " " + userLevel + " " + userPlaycount + " " + userHours + " " + userMedalCount + " " + userPeakRank
	client.Reply(returnString, f)
	return
}

func IntToCommaString(n int) string {
	// Convert the integer to a string
	str := strconv.Itoa(n)

	// Find the length of the string
	length := len(str)

	// If the string length is less than or equal to 3, return it as is
	if length <= 3 {
		return str
	}

	// Create a slice to hold the characters of the result
	var result []byte

	// Calculate the position of the first comma
	firstComma := length % 3
	if firstComma == 0 {
		firstComma = 3
	}

	// Append the first segment before the first comma
	result = append(result, str[:firstComma]...)

	// Iterate over the rest of the string, adding commas
	for i := firstComma; i < length; i += 3 {
		result = append(result, ',')
		result = append(result, str[i:i+3]...)
	}

	return string(result)
}

func FloatToCommaString(n float64) string {
	// Format the float to a string with 2 decimal places
	str := fmt.Sprintf("%.2f", n)

	// Split the string into integer and fractional parts
	parts := strings.Split(str, ".")
	integerPart := parts[0]
	fractionalPart := parts[1]

	// Handle the integer part (adding commas)
	var result []byte
	length := len(integerPart)
	firstComma := length % 3
	if firstComma == 0 {
		firstComma = 3
	}
	result = append(result, integerPart[:firstComma]...)
	for i := firstComma; i < length; i += 3 {
		result = append(result, ',')
		result = append(result, integerPart[i:i+3]...)
	}

	// Combine the integer part with the fractional part
	return string(result) + "." + fractionalPart
}
