package aprsHelper

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type APRSUserClient struct {
	callSign     string
	ssid         int
	password     int
	messageQueue MessageQueue
}

func ExtractCommand(message string) string {
	// Remove the header (everything before and including the first space)
	parts := strings.SplitN(message, " :", 2)
	if len(parts) < 2 {
		return ""
	}
	messageBody := parts[1]

	// Remove the footer (everything after and including the '{')
	messageBody = strings.SplitN(messageBody, "{", 2)[0]

	if strings.HasPrefix(messageBody, "!") {
		messageBody = strings.TrimPrefix(messageBody, "!")
	}

	// Return the cleaned-up message
	return strings.TrimSpace(messageBody)
}

func ExtractArgs(message string) ([]string, error) {
	// Remove the leading '!'
	message = strings.TrimPrefix(message, "!")

	// Split the command and its arguments by spaces
	args := strings.Fields(message)

	// Ensure that there's at least a command present
	if len(args) == 0 {
		return nil, fmt.Errorf("no command found in the message")
	}

	// Return the arguments slice
	return args[1:], nil
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

func (client APRSUserClient) SendAck(f aprs.Frame) {
	messageNum, _ := extractMessageNumber(f.Text)
	personWhoMessagedMe := ExtractAuthor(f)
	botStation := aprs.Addr{
		Call: client.callSign,
		SSID: client.ssid,
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
	sendMessageFrame(ack)
}

func EnsureLength(input string) string {
	if len(input) >= 9 {
		return input[:9] // Truncate if longer than 9 characters
	}
	return input + spaces(9-len(input)) // Append spaces to reach 9 characters
}

// spaces returns a string of the specified length consisting of spaces.
func spaces(n int) string {
	return " " + string(make([]rune, n-1)) // Create a string with n spaces
}

func (client APRSUserClient) GenerateMessageReplyFrame(messageContent string, f aprs.Frame) aprs.Frame {
	personWhoMessagedMe := ExtractAuthor(f)
	botStation := aprs.Addr{
		Call: client.callSign,
		SSID: client.ssid,
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

func ExtractAuthor(frame aprs.Frame) string {
	var author = frame.Src.Call + "-" + strconv.Itoa(frame.Src.SSID)
	return author
}

func sendMessageFrame(f aprs.Frame) {
	err := f.SendIS("tcp://rotate.aprs.net:14580", 24233)
	if err != nil {
		fmt.Println("Failed to send message to APRSIS: " + err.Error())
		return
	}
}

func extractSSIDFromCallSSID(input string) (string, int) {
	var split = strings.Split(input, "-")
	var callSign = split[0]
	var ssid, err = strconv.Atoi(split[1])
	if err != nil {
		panic("Failed to pull callsign and ssid")
	}
	return callSign, ssid
}

func InitAPRSClient(callandSSID string, APRSPassword int) APRSUserClient {
	var call, ssid = extractSSIDFromCallSSID(callandSSID)
	var messageQueue = NewMessageQueue()
	return APRSUserClient{callSign: call, ssid: ssid, password: APRSPassword, messageQueue: messageQueue}
}

func (client APRSUserClient) AprsTextReply(text string, f aprs.Frame) {
	if len(text) <= 67 {
		// instead of directly sending the messages, lets have a queueing system that the messages get added to.
		// in this queue, we can listen for acks and all. We can also then monitor the queue to see how many messages we
		// process, as well as rate-limit ourselves.
		client.messageQueue.Push(client.GenerateMessageReplyFrame(text, f))
		//sendMessageFrame(client.GenerateMessageReplyFrame(text, f))
	} else {
		// split the frame text into several packets and send each of them
		var messages = splitStringByLength(text, 66)
		fmt.Println("message split")
		fmt.Println(messages)
		for _, message := range messages {
			//fmt.Println("sending message", i, message)
			sendMessageFrame(client.GenerateMessageReplyFrame(message, f))
			time.Sleep(3 * time.Second)
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
