package general

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"simpleAPRSbot-go/aprsHelper"
	"strconv"
	"strings"
)

func Ping(args []string, f aprs.Frame) {
	replyText := "Pong!"
	if len(args) != 0 {
		replyText = strings.Join(args, " ")
	}
	// reply, I threw in a magic int here to make the message unique always
	replyMessageFrame := aprsHelper.GenerateMessageReplyFrame(replyText+" "+strconv.Itoa(rand.IntN(1000)), f)
	fmt.Println(replyMessageFrame)
	//time.Sleep(2 * time.Second)
	aprsHelper.SendMessage(replyMessageFrame)
	return
}
