package general

import (
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"simpleAPRSbot-go/helpers/aprsHelper"
)

func Flip(args []string, f aprs.Frame, client *aprsHelper.APRSUserClient) {
	// heads or tails
	var options = [2]string{"Heads", "Tails"}
	var decideInt = rand.IntN(1000)
	result := options[decideInt%2]
	client.AprsTextReply(result, f)
}
