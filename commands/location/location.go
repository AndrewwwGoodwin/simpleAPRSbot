package location

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"simpleAPRSbot-go/helpers/APRS"
)

func Location(args []string, f aprs.Frame, client *APRS.UserClient) {
	// this command gets the user's last seen location, and returns their current zip code.
	// step 1: get the user's last location

	var callerCallsign string
	switch len(args) {
	case 0:
		callerCallsign = APRS.ExtractAuthor(f)
	case 1:
		callerCallsign = args[0]
	default:
		client.AprsTextReply("Too many args", f)
		return
	}

	// now that we know who to look for, lets find them!
	// ring up the APRSFi API
	var wrapper = client.ApiClients.APRSFi
	var data, err = wrapper.GetLocation(callerCallsign)
	if err != nil {
		fmt.Println(err)
		client.AprsTextReply(err.Error(), f)
		return
	}
	var locationData = data.Entries[0]

	// for now lets just return their lat/long
	// in the future I want to send this lat/long off to some geocoding api
	// which will then return a lot more detailed information such as city,state zipcode, county, coordinates,
	// how long ago, via aprs.fi
	client.AprsTextReply(locationData.Lat+" "+locationData.Lng+" via aprs.fi", f)
	return
}
