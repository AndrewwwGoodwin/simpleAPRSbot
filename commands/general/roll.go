package general

import (
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"simpleAPRSbot-go/aprsHelper"
	"strconv"
)

func Roll(args []string, f aprs.Frame) {
	aprsHelper.AprsTextReply(strconv.Itoa(rand.IntN(100)), f)
	return
}
