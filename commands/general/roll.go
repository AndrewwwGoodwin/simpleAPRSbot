package general

import (
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"simpleAPRSbot-go/helpers/aprsHelper"
	"strconv"
)

func Roll(args []string, f aprs.Frame) {
	// take in 1 arg that would allow for setting max roll
	// otherwise default to rolling out of 100
	var maxRoll = 100
	if len(args) > 0 {
		var err error
		maxRoll, err = strconv.Atoi(args[0])
		if err != nil || maxRoll <= 0 {
			// Handle error or invalid maxRoll value, default to 100
			maxRoll = 100
		}
	}
	aprsHelper.AprsTextReply(strconv.Itoa(rand.IntN(maxRoll)), f)
	return
}
