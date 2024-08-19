package aprsHelper

import "github.com/ebarkie/aprs"

type Message struct {
	Message                 aprs.Frame
	ReceivedAcknowledgement bool
}

type MessageQueue struct {
	queue []aprs.Frame
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{}
}

func (mq *MessageQueue) Push(f aprs.Frame) {
	mq.queue = append(mq.queue, f)
}
func (mq *MessageQueue) Pop() aprs.Frame {
	if len(mq.queue) == 0 {
		return aprs.Frame{}
	} else {
		f := mq.queue[0]
		mq.queue = mq.queue[1:]
		return f
	}
}

func (mq *MessageQueue) Clear() {
	mq.queue = []aprs.Frame{}
}
