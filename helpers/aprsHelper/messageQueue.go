package aprsHelper

import (
	"fmt"
	"github.com/ebarkie/aprs"
)

type Message struct {
	Message                 aprs.Frame
	ReceivedAcknowledgement bool
}

type MessageQueue struct {
	queue []Message
}

func NewMessageQueue() MessageQueue {
	var newQueue = MessageQueue{make([]Message, 0)}
	return newQueue
}

func (mq MessageQueue) Push(f aprs.Frame) {
	var m = Message{
		Message:                 f,
		ReceivedAcknowledgement: false,
	}
	fmt.Println(mq.queue, m)
	mq.queue = append(mq.queue, m)
	fmt.Println(mq.queue)
}

func (mq MessageQueue) Pop() aprs.Frame {
	if len(mq.queue) == 0 {
		return aprs.Frame{}
	} else {
		f := mq.queue[0].Message
		mq.queue = mq.queue[1:]
		return f
	}
}
