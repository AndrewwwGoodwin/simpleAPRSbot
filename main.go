package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ebarkie/aprs"
	"log"
	"os"
	"os/signal"
	"simpleAPRSbot-go/aprsHelper"
	"simpleAPRSbot-go/commands/general"
	"strings"
)

type CommandFunc func(args []string, f aprs.Frame)

type CommandFuncAPRSFi func(args []string, f aprs.Frame, aprsFiKey string)

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

var commandRegistryAPRSFI = map[string]CommandFuncAPRSFi{
	"loc":      general.Location,
	"location": general.Location,
}

var aprsCALL = flag.String("APRS_CALL", "N0CALL-0", "N0CALL-0")
var aprsPass = flag.Int("APRS_PASS", 000000, "00000")
var AprsFiAPIKey = flag.String("APRS_FI_API_KEY", "", "")

func main() {
	// waits for termination so everything shuts down nicely
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Shutting down")
		os.Exit(0)
	}()
	flag.Parse()
	fmt.Println(aprsCALL, aprsPass)

	log.Println("Receiving")
	for {
		ctx := context.Background()
		fc := aprs.RecvIS(ctx, "rotate.aprs.net:14580", aprs.Addr{Call: *aprsCALL}, *aprsPass, "g/"+*aprsCALL)
		for f := range fc {
			fmt.Println("")
			fmt.Println(f)
			// dont ack acks
			if strings.HasPrefix(f.Text, ":"+aprsHelper.EnsureLength(*aprsCALL)+":ack") {
				continue
			}
			aprsHelper.SendAck(f)
			if strings.HasPrefix(f.Text, ":"+aprsHelper.EnsureLength(*aprsCALL)+":!") {
				//strip the prefix
				commandName := strings.ToLower(strings.Split(aprsHelper.ExtractCommand(f.Text), " ")[0])
				commandArgs, _ := aprsHelper.ExtractArgs(aprsHelper.ExtractCommand(f.Text))
				if _, exists := commandRegistry[commandName]; exists {
					// The command exists, you can execute it
					handleCommand(commandName, commandArgs, f)
				} else if _, exists := commandRegistryAPRSFI[commandName]; exists {
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
	} else if commandFuncAPRSFi, existsAprs := commandRegistryAPRSFI[strings.ToLower(commandName)]; existsAprs {
		commandFuncAPRSFi(commandArgs, f, *AprsFiAPIKey)
	} else {
		fmt.Println("Unknown command:", commandName)
	}
}
