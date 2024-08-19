package general

import (
	"github.com/ebarkie/aprs"
	"simpleAPRSbot-go/helpers/aprsHelper"
	"strings"
)

func Ping(args []string, f aprs.Frame, client *aprsHelper.APRSUserClient) {
	replyText := "Pong!"
	if len(args) != 0 {
		replyText = strings.Join(args, " ")
	}
	// reply, I threw in a magic int here to make the message unique always
	//time.Sleep(2 * time.Second)
	client.AprsTextReply(replyText, f)
	return
}
