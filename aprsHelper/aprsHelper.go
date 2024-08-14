package aprsHelper

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"strings"
)

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

func SendAck(f aprs.Frame) {
	messageNum, _ := extractMessageNumber(f.Text)
	personWhoMessagedMe, _ := ExtractAuthor(f.String())
	botStation := aprs.Addr{
		Call: "KQ4NRT",
		SSID: 6,
	}
	botToCall := aprs.Addr{
		Call: "APZ727",
	}
	//botPath := aprs.Path{
	//	aprs.Addr{
	//		Call: "TCPIP*",
	//	},
	//}
	// ack the message
	ack := aprs.Frame{
		Dst: botToCall,
		Src: botStation,
		//Path: botPath,
		Text: ":" + EnsureLength(personWhoMessagedMe) + ":" + "ack" + messageNum,
	}
	fmt.Println(ack)
	SendMessage(ack)
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

func GenerateMessageReplyFrame(messageContent string, f aprs.Frame) aprs.Frame {
	personWhoMessagedMe, _ := ExtractAuthor(f.String())
	botStation := aprs.Addr{
		Call: "KQ4NRT",
		SSID: 6,
	}
	botToCall := aprs.Addr{
		Call: "APZ727",
	}
	messageFrame := aprs.Frame{
		Dst: botToCall,
		Src: botStation,
		//Path: botPath,
		Text: ":" + EnsureLength(personWhoMessagedMe) + ":" + messageContent,
	}
	return messageFrame
}

func ExtractAuthor(frame string) (string, error) {
	// Find the position of the '>' symbol that separates the author from the destination
	greaterThanIndex := strings.Index(frame, ">")
	if greaterThanIndex == -1 {
		return "", fmt.Errorf("no '>' found in the frame")
	}

	// Extract the author, which is everything before the '>'
	author := frame[:greaterThanIndex]

	// Ensure that there is an author
	if len(author) == 0 {
		return "", fmt.Errorf("no author found before '>'")
	}

	return author, nil
}

func SendMessage(f aprs.Frame) {
	frameLength := len(f.Text) - 11

	// max message length is 67 characters, so if its longer we need to split
	if frameLength <= 67 {
		err := f.SendIS("tcp://rotate.aprs.net:14580", 24233)
		if err != nil {
			fmt.Println("Failed to send message to APRSIS: " + err.Error())
			return
		}
	} else {
		// split the frame text into several packets and send each of them
		AprsTextReply("Reply message too long.", f)
	}
}

func AprsTextReply(text string, f aprs.Frame) {
	SendMessage(GenerateMessageReplyFrame(text, f))
	return
}
