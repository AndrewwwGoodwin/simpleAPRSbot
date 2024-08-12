package general

import (
	"github.com/ebarkie/aprs"
	"osuAPRS/aprsHelper"
	"time"
)

//todo: add timezone support via args

func Time(args []string, f aprs.Frame) {
	aprsHelper.AprsTextReply(time.Now().UTC().Format("02 Jan 06 15:04:05 MST"), f)
	return
}
