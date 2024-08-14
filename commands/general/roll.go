package general

import (
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"simpleAPRSbot-go/aprsHelper"
	"strconv"
)

func Roll(args []string, f aprs.Frame) {
	// take in 1 arg that would allow for setting max roll
	// otherwise default to rolling out of 100

	aprsHelper.AprsTextReply(strconv.Itoa(rand.IntN(100)), f)
	return
}
