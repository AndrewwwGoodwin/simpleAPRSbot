package aprsHelper

import (
	"github.com/ebarkie/aprs"
)

type Message struct {
	Message                 *aprs.Frame
	ReceivedAcknowledgement bool
}

type MessageQueue struct {
	Queue []Message
}

func NewMessageQueue() *MessageQueue {
	var newQueue = MessageQueue{make([]Message, 0)}
	return &newQueue
}

func (mq *MessageQueue) Push(f aprs.Frame) {
	var m = Message{
		Message:                 &f,
		ReceivedAcknowledgement: false,
	}
	mq.Queue = append(mq.Queue, m)
}

func (mq *MessageQueue) Pop() aprs.Frame {
	if len(mq.Queue) == 0 {
		return aprs.Frame{}
	} else {
		f := mq.Queue[0].Message
		mq.Queue = mq.Queue[1:]
		return *f
	}
}
