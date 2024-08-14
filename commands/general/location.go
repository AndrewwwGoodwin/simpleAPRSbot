package general

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"simpleAPRSbot-go/aprsFiWrapper"
	"simpleAPRSbot-go/aprsHelper"
	"strconv"
)

func Location(args []string, f aprs.Frame, apiKey string) {
	// this command gets the user's last seen location, and returns their current zip code.
	// step 1: get the user's last location

	var callerCallsign string
	switch len(args) {
	case 0:
		callerCallsign, _ = aprsHelper.ExtractAuthor(f.String())
	case 1:
		callerCallsign = args[0]
	default:
		aprsHelper.AprsTextReply("Too many args", f)
		return
	}

	// now that we know who to look for, lets find them!
	// ring up the APRSFi API
	var wrapper = aprsFiWrapper.NewAprsFiWrapper(apiKey)
	var data, err = wrapper.GetLocation(callerCallsign)
	if err != nil {
		fmt.Println(err)
		aprsHelper.AprsTextReply(err.Error(), f)
		return
	}
	var locationData = data.Entries[0]

	// for now lets just return their lat/long
	// in the future I want to send this lat/long off to some geocoding api
	// which will then return a lot more detailed information such as city,state zipcode, county, coordinates,
	// how long ago, via aprs.fi
	aprsHelper.AprsTextReply(locationData.Lat+" "+locationData.Lng+" via aprs.fi"+" "+strconv.Itoa(rand.IntN(999)), f)
	return
}
