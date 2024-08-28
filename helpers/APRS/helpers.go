package APRS

import (
	"errors"
	"fmt"
	"github.com/ebarkie/aprs"
	"strings"
)

func GetAuthor(frame aprs.Frame) string {
	var author = frame.Src.String()
	return author
}

func SendMessageFrame(f aprs.Frame) {
	err := f.SendIS("tcp://rotate.aprs.net:14580", 24233)
	if err != nil {
		fmt.Println("Failed to send message to APRSIS: " + err.Error())
		return
	}
}

type Command struct {
	Name      string
	Arguments []string
}

func GetCommand(message string) (*Command, error) {
	// Remove the header (everything before and including the first space)
	parts := strings.SplitN(message, " :", 2)
	if len(parts) < 2 {
		return nil, errors.New("Invalid Packet: " + message)
	}
	messageBody := parts[1]

	// Remove the footer (everything after and including the '{')
	messageBody = strings.SplitN(messageBody, "{", 2)[0]

	if strings.HasPrefix(messageBody, "!") {
		messageBody = strings.TrimPrefix(messageBody, "!")
	}

	var commandAndArgs = strings.Split(strings.TrimSpace(messageBody), " ")

	returnData := Command{
		Name:      strings.ToLower(commandAndArgs[0]),
		Arguments: commandAndArgs[1:],
	}

	return &returnData, nil
}

func EnsureLength(input string) string {
	if len(input) >= 9 {
		return input[:9] // Truncate if longer than 9 characters
	}
	return input + strings.Repeat(" ", 9-len(input)) // Append spaces to reach 9 characters
}
