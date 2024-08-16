package general

import (
	"github.com/ebarkie/aprs"
	"regexp"
	"simpleAPRSbot-go/aprsHelper"
	"strings"
)

func CW(args []string, f aprs.Frame) {
	var toEncode = strings.Join(args, " ")
	// use regex to ensure that the input string can be converted to morse
	var isValid = regexp.MustCompile(`^[a-zA-Z0-9 .,?'!/()]+$`).MatchString(toEncode)

	if isValid {
		// so if the text is valid, we need to generate a string representing morse code for the input
		var morseString = generateMorse(toEncode)
		aprsHelper.AprsTextReply(morseString, f)
		return
	} else {
		aprsHelper.AprsTextReply("Invalid Characters found in input. a-z A-Z 0-9 only.", f)
	}
}

func generateMorse(toEncode string) string {
	// this map is our lookup table for each character.
	// thank you chatGPT for typing that for me
	lookupTable := map[string]string{
		"A": ".-", "B": "-...", "C": "-.-.", "D": "-..", "E": ".", "F": "..-.",
		"G": "--.", "H": "....", "I": "..", "J": ".---", "K": "-.-", "L": ".-..",
		"M": "--", "N": "-.", "O": "---", "P": ".--.", "Q": "--.-", "R": ".-.",
		"S": "...", "T": "-", "U": "..-", "V": "...-", "W": ".--", "X": "-..-",
		"Y": "-.--", "Z": "--..",

		"0": "-----", "1": ".----", "2": "..---", "3": "...--", "4": "....-",
		"5": ".....", "6": "-....", "7": "--...", "8": "---..", "9": "----.",

		".": ".-.-.-", ",": "--..--", "?": "..--..", "'": ".----.", "!": "-.-.--",
		"/": "-..-.", "(": "-.--.", ")": "-.--.-", " ": "/",
	}
	var returnString string
	for _, element := range strings.ToUpper(toEncode) {
		morseCode, exists := lookupTable[string(element)]
		if exists {
			returnString += morseCode + " "
		}
	}
	return returnString
}
