package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ebarkie/aprs"
	"log"
	"os"
	"os/signal"
	"simpleAPRSbot-go/api"
	"simpleAPRSbot-go/aprsHelper"
	"simpleAPRSbot-go/commands/general"
	"simpleAPRSbot-go/commands/location"
	"strings"
)

type CommandFunc func(args []string, f aprs.Frame)

type CommandFuncAPIKeys func(args []string, f aprs.Frame, aprsFiKey api.Keys)

var commandRegistry = map[string]CommandFunc{
	"ping":     general.Ping,
	"p":        general.Ping,
	"time":     general.Time,
	"t":        general.Time,
	"flip":     general.Flip,
	"coinflip": general.Flip,
	"roll":     general.Roll,
	"r":        general.Roll,
	"cw":       general.CW,
	"calc":     general.CalculateCommand,
	"eval":     general.CalculateCommand,
}

var commandRegistryAPRSFI = map[string]CommandFuncAPIKeys{
	"loc":      location.Location,
	"location": location.Location,
	"w":        location.Weather,
	"weather":  location.Weather,
}

var aprsCALL = flag.String("APRS_CALL", "N0CALL-0", "N0CALL-0")
var aprsPass = flag.Int("APRS_PASS", 000000, "00000")

func main() {
	// waits for termination so everything shuts down nicely
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Shutting down")
		os.Exit(0)
	}()
	// load up our flags
	APRSFIkey := flag.String("APRS_FI_API_KEY", "", "APRS FI API Key")
	OpenWeatherMapkey := flag.String("OWM_FI_API_KEY", "", "OpenWeatherMap API Key")
	flag.Parse()

	//shove them in an object we can ship around the program
	var APIKeyObj api.Keys
	APIKeyObj.APRSFIkey = *APRSFIkey
	APIKeyObj.OpenWeatherMapkey = *OpenWeatherMapkey

	log.Println("Receiving")
	for {
		ctx := context.Background()
		fc := aprs.RecvIS(ctx, "rotate.aprs.net:14580", aprs.Addr{Call: *aprsCALL}, *aprsPass, "g/"+*aprsCALL)
		for f := range fc {
			fmt.Println("")
			fmt.Println(f)
			if strings.HasPrefix(f.Text, ":"+aprsHelper.EnsureLength(*aprsCALL)+":!") {
				aprsHelper.SendAck(f)
				//strip the prefix
				commandName := strings.ToLower(strings.Split(aprsHelper.ExtractCommand(f.Text), " ")[0])
				commandArgs, _ := aprsHelper.ExtractArgs(aprsHelper.ExtractCommand(f.Text))
				if commandFunc, exists := commandRegistry[commandName]; exists {
					commandFunc(commandArgs, f) // Call the corresponding function
				} else if commandFuncAPRSFi, existsAprs := commandRegistryAPRSFI[strings.ToLower(commandName)]; existsAprs {
					commandFuncAPRSFi(commandArgs, f, APIKeyObj)
				} else {
					fmt.Println("Unknown command:", commandName)
				}
			} else {
				// dont ack acks
				if strings.HasPrefix(f.Text, ":"+aprsHelper.EnsureLength(*aprsCALL)+":ack") {
					continue
					// dont ack messages not sent to us
				} else if !strings.HasPrefix(f.Text, ":"+aprsHelper.EnsureLength(*aprsCALL)+":") {
					continue
				} else {
					// if we make it through all that, finally ack the message
					aprsHelper.SendAck(f)
				}
			}
		}
	}
}
