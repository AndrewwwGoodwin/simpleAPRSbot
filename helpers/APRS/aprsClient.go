package APRS

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"simpleAPRSbot-go/helpers/api"
	"strconv"
	"strings"
	"unicode/utf8"
)

type UserClient struct {
	CallSign        string
	APRSSsid        int
	APRSCallAndSSID string
	APRSPassword    int
	MessageQueue    *MessageQueue
	ApiClients      api.Clients
}

func (client UserClient) SendAck(f aprs.Frame) {
	messageNum, _ := extractMessageNumber(f.Text)
	personWhoMessagedMe := GetAuthor(f)
	botStation := aprs.Addr{
		Call: client.CallSign,
		SSID: client.APRSSsid,
	}
	botToCall := aprs.Addr{
		Call: "APZ727",
	}
	var messageText = "ack" + messageNum
	// ack the message
	ack := aprs.Frame{
		Dst: botToCall,
		Src: botStation,
		//Path: botPath,
		Text: ":" + EnsureLength(personWhoMessagedMe) + ":" + messageText,
	}
	fmt.Println(ack)
	client.MessageQueue.Push(ack)
}

func (client UserClient) GenerateMessageReplyFrame(messageContent string, f aprs.Frame) aprs.Frame {
	personWhoMessagedMe := GetAuthor(f)
	botStation := aprs.Addr{
		Call: client.CallSign,
		SSID: client.APRSSsid,
	}
	botToCall := aprs.Addr{
		Call: "APZ727",
	}
	messageFrame := aprs.Frame{
		Dst: botToCall,
		Src: botStation,
		//Path: botPath,
		Text: ":" + EnsureLength(personWhoMessagedMe) + ":" + messageContent + "{" + strconv.Itoa(rand.IntN(999)),
	}
	fmt.Println(messageFrame.String())
	return messageFrame
}

func extractSSIDFromCallSSID(input string) (string, int) {
	var split = strings.Split(input, "-")
	var callSign = split[0]
	var ssid, err = strconv.Atoi(split[1])
	if err != nil {
		panic("Failed to pull callsign and APRSSsid")
	}
	return callSign, ssid
}

func InitAPRSClient(callandSSID string, APRSPassword int, apiClients api.Clients) *UserClient {
	var call, ssid = extractSSIDFromCallSSID(callandSSID)
	var messageQueue = NewMessageQueue()
	return &UserClient{CallSign: call, APRSCallAndSSID: callandSSID, APRSSsid: ssid, APRSPassword: APRSPassword, MessageQueue: messageQueue, ApiClients: apiClients}
}

func (client UserClient) Reply(text string, f aprs.Frame) {
	if len(text) <= 67 {
		// instead of directly sending the messages, lets have a queueing system that the messages get added to.
		// in this Queue, we can listen for acks and all. We can also then monitor the Queue to see how many messages we
		// process, as well as rate-limit ourselves.
		client.MessageQueue.Push(client.GenerateMessageReplyFrame(text, f))
		//SendMessageFrame(client.GenerateMessageReplyFrame(text, f))
	} else {
		// split the frame text into several packets and send each of them
		var messages = splitStringByLength(text, 66)
		fmt.Println("message split")
		fmt.Println(messages)
		var packets = make([]aprs.Frame, len(messages))
		for _, message := range messages {
			//fmt.Println("sending message", i, message)
			packets = append(packets, client.GenerateMessageReplyFrame(message, f))
		}
		for _, packet := range packets {
			client.MessageQueue.Push(packet)
		}
		return
	}
}

func splitStringByLength(s string, maxLength int) []string {
	var result []string
	for len(s) > 0 {
		if len(s) <= maxLength {
			result = append(result, s)
			break
		}

		// Try to find the last space within the maxLength boundary
		cutIndex := strings.LastIndex(s[:maxLength], " ")
		if cutIndex == -1 {
			// No space found, fall back to cutting at maxLength
			cutIndex = maxLength
			// Ensure we're not cutting a multibyte character
			for !utf8.ValidString(s[:cutIndex]) {
				cutIndex--
			}
		}

		// Append the section and trim the string
		result = append(result, s[:cutIndex])
		s = strings.TrimSpace(s[cutIndex:])
	}
	return result
}

func extractMessageNumber(message string) (string, error) {
	// Find the last '{' in the message
	lastBraceIndex := strings.LastIndex(message, "{")
	if lastBraceIndex == -1 {
		return "", fmt.Errorf("no '{' found in the message")
	}

	// Extract everything after the last '{'
	messageNumber := message[lastBraceIndex+1:]

	// Ensure that there is a message number after '{'
	if len(messageNumber) == 0 {
		return "", fmt.Errorf("no message number found after '{'")
	}

	return messageNumber, nil
}
