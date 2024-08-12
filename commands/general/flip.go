package general

import (
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"simpleAPRSbot-go/aprsHelper"
	"strconv"
)

func Flip(args []string, f aprs.Frame) {
	// heads or tails
	var options = [2]string{"Heads", "Tails"}
	var decideInt = rand.IntN(1000)
	result := options[decideInt%2]
	aprsHelper.AprsTextReply(result+" "+strconv.Itoa(decideInt), f)
}
