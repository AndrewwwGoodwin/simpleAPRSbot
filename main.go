package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ebarkie/aprs"
	"log"
	"simpleAPRSbot-go/aprsHelper"
	"simpleAPRSbot-go/commands/general"
	"strings"
)

type CommandFunc func(args []string, f aprs.Frame)

var commandRegistry = map[string]CommandFunc{
	"ping":     general.Ping,
	"p":        general.Ping,
	"time":     general.Time,
	"t":        general.Time,
	"flip":     general.Flip,
	"coinflip": general.Flip,
	"roll":     general.Roll,
	"r":        general.Roll,
}

func main() {
	var aprsCALL = flag.String("APRS_CALL", "N0CALL-0", "N0CALL-0")
	var aprsPass = flag.Int("APRS_PASS", 000000, "00000")
	flag.Parse()
	fmt.Println(aprsCALL, aprsPass)

	log.Println("Receiving")
	for {
		ctx := context.Background()
		fc := aprs.RecvIS(ctx, "rotate.aprs.net:14580", aprs.Addr{Call: *aprsCALL}, *aprsPass, "e/"+*aprsCALL)
		for f := range fc {
			fmt.Println("")
			fmt.Println(f)

			aprsHelper.SendAck(f)
			if strings.HasPrefix(f.Text, ":"+aprsHelper.EnsureLength(*aprsCALL)+":!") {
				//strip the prefix
				commandName := strings.ToLower(strings.Split(aprsHelper.ExtractCommand(f.Text), " ")[0])
				commandArgs, _ := aprsHelper.ExtractArgs(aprsHelper.ExtractCommand(f.Text))
				if _, exists := commandRegistry[commandName]; exists {
					// The command exists, you can execute it
					handleCommand(commandName, commandArgs, f)
				} else {
					fmt.Println("Unknown command:", commandName)
				}
			}
		}
	}
}

func handleCommand(commandName string, commandArgs []string, f aprs.Frame) {
	if commandFunc, exists := commandRegistry[commandName]; exists {
		commandFunc(commandArgs, f) // Call the corresponding function
	} else {
		fmt.Println("Unknown command:", commandName)
	}
}
