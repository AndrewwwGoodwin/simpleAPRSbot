package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ebarkie/aprs"
	"log"
	"os"
	"os/signal"
	"simpleAPRSbot-go/commands/general"
	"simpleAPRSbot-go/commands/location"
	osuCommands "simpleAPRSbot-go/commands/osu"
	"simpleAPRSbot-go/helpers/api"
	"simpleAPRSbot-go/helpers/aprsHelper"
	"strings"
	"time"
)

type CommandFunc func(args []string, f aprs.Frame, client *aprsHelper.APRSUserClient)

type CommandFuncAPIKeys func(args []string, f aprs.Frame, aprsFiKey api.Keys, client *aprsHelper.APRSUserClient)

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

var commandRegistryAPIKeysRequired = map[string]CommandFuncAPIKeys{
	"loc":      location.Location,
	"location": location.Location,
	"w":        location.Weather,
	"weather":  location.Weather,
	"osu":      osuCommands.Osu,
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
	OpenWeatherMapKey := flag.String("OWM_FI_API_KEY", "", "OpenWeatherMap API Key")
	osuClientID := flag.Int("OSU_CLIENT_ID", 0, "OSU Client ID")
	osuClientSecret := flag.String("OSU_CLIENT_SECRET", "", "OSU Client Secret")
	flag.Parse()

	//shove them in an object we can ship around the program
	var APIKeyObj api.Keys
	APIKeyObj.APRSFIkey = *APRSFIkey
	APIKeyObj.OpenWeatherMapkey = *OpenWeatherMapKey
	APIKeyObj.OsuClientID = *osuClientID
	APIKeyObj.OsuClientSecret = *osuClientSecret

	// we also need to create an instance of APRSUserClient, so we can reply to messages
	var client = aprsHelper.InitAPRSClient(*aprsCALL, *aprsPass)

	// crank up the queue processor
	go queueProcessor(client)

	// and start listening for commands
	go commandListener(client, APIKeyObj)
	log.Println("Receiving")

	select {}
}

func queueProcessor(client *aprsHelper.APRSUserClient) {
	for {
		if len(client.MessageQueue.Queue) <= 0 {
			continue
		} else {
			aprsHelper.SendMessageFrame(client.MessageQueue.Pop())
			// this globally lets us only send a message every x secs. can be turned up or down based on load
			time.Sleep(1 * time.Second)
		}
	}
}

func commandListener(client *aprsHelper.APRSUserClient, apiKeyObj api.Keys) {
	for {
		ctx := context.Background()
		fc := aprs.RecvIS(ctx, "rotate.aprs.net:14580", aprs.Addr{Call: *aprsCALL}, *aprsPass, "g/"+*aprsCALL)
		for f := range fc {
			fmt.Println("")
			fmt.Println(f)
			if strings.HasPrefix(f.Text, ":"+aprsHelper.EnsureLength(*aprsCALL)+":!") {
				client.SendAck(f)
				//strip the prefix
				commandName := strings.ToLower(strings.Split(aprsHelper.ExtractCommand(f.Text), " ")[0])
				commandArgs, _ := aprsHelper.ExtractArgs(aprsHelper.ExtractCommand(f.Text))
				if commandFunc, exists := commandRegistry[commandName]; exists {
					commandFunc(commandArgs, f, client) // Call the corresponding function
				} else if commandFuncAPRSFi, existsAprs := commandRegistryAPIKeysRequired[strings.ToLower(commandName)]; existsAprs {
					commandFuncAPRSFi(commandArgs, f, apiKeyObj, client)
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
					client.SendAck(f)
				}
			}
		}
	}
}
