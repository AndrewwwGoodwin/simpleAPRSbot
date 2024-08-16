package general

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"simpleAPRSbot-go/aprsHelper"
	"strings"
	"time"
)

// add timezone support via args
var timezoneAbbrs = map[string]string{
	"UTC": "UTC",
	"EST": "America/New_York",
	"EDT": "America/New_York",
	"CST": "America/Chicago",
	"CDT": "America/Chicago",
	"MST": "America/Denver",
	"MDT": "America/Denver",
	"PST": "America/Los_Angeles",
	"PDT": "America/Los_Angeles",
}

func Time(args []string, f aprs.Frame) {
	// Default timezone to UTC
	location := "UTC"

	// If arguments are provided, use the first one as the timezone
	if len(args) > 0 {
		arg := strings.ToUpper(args[0])
		if mappedLocation, exists := timezoneAbbrs[arg]; exists {
			location = mappedLocation
		} else {
			// Handle error if the timezone abbreviation is invalid
			aprsHelper.AprsTextReply(fmt.Sprintf("Error: Invalid timezone abbreviation %s", arg), f)
			return
		}
	}

	// Parse the timezone
	loc, err := time.LoadLocation(location)
	if err != nil {
		// Handle error if the timezone is invalid
		aprsHelper.AprsTextReply(fmt.Sprintf("Error: Invalid timezone %s", location), f)
		return
	}

	// Get the current time in the specified location
	currentTime := time.Now().In(loc)
	formattedTime := currentTime.Format("02 Jan 06 15:04:05 MST")

	// Send the formatted time as a reply
	aprsHelper.AprsTextReply(formattedTime, f)
}
