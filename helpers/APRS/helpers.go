package APRS

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"strconv"
	"strings"
)

func ExtractAuthor(frame aprs.Frame) string {
	var author = frame.Src.Call + "-" + strconv.Itoa(frame.Src.SSID)
	return author
}

func SendMessageFrame(f aprs.Frame) {
	err := f.SendIS("tcp://rotate.aprs.net:14580", 24233)
	if err != nil {
		fmt.Println("Failed to send message to APRSIS: " + err.Error())
		return
	}
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

func EnsureLength(input string) string {
	if len(input) >= 9 {
		return input[:9] // Truncate if longer than 9 characters
	}
	return input + spaces(9-len(input)) // Append spaces to reach 9 characters
}
