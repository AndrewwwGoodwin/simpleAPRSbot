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
	"loc":      location.Location,
	"location": location.Location,
	"w":        location.Weather,
	"weather":  location.Weather,
	"osu":      osuCommands.Osu,
}

var APRSClient *aprsHelper.APRSUserClient

func init() {
	// load in all our flags
	var aprsCALL = flag.String("APRS_CALL", "N0CALL-0", "N0CALL-0")
	var aprsPass = flag.Int("APRS_PASS", 000000, "00000")
	var APRSFIkey = flag.String("APRS_FI_API_KEY", "", "APRS FI API Key")
	var OpenWeatherMapKey = flag.String("OWM_FI_API_KEY", "", "OpenWeatherMap API Key")
	var osuClientID = flag.Int("OSU_CLIENT_ID", 0, "OSU Client ID")
	var osuClientSecret = flag.String("OSU_CLIENT_SECRET", "", "OSU Client Secret")
	flag.Parse()
	// build our objects
	APIClients := api.InitializeAPIClients(api.Keys{
		OsuClientSecret:   *osuClientSecret,
		OsuClientID:       *osuClientID,
		OpenWeatherMapKey: *OpenWeatherMapKey,
		APRSFIkey:         *APRSFIkey,
	})

	APRSClient = aprsHelper.InitAPRSClient(*aprsCALL, *aprsPass, APIClients)
}

func main() {
	// waits for termination so everything shuts down nicely
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Shutting down")
		os.Exit(0)
	}()

	// crank up the queue processor
	go queueProcessor(APRSClient)

	// and start listening for commands
	go commandListener(APRSClient)
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

func commandListener(client *aprsHelper.APRSUserClient) {
	for {
		ctx := context.Background()
		fc := aprs.RecvIS(ctx, "rotate.aprs.net:14580", aprs.Addr{Call: client.CallSign, SSID: client.APRSSsid}, client.APRSPassword, "g/"+client.APRSCallAndSSID)
		for receivedMessageFrame := range fc {
			fmt.Println("")
			fmt.Println(receivedMessageFrame)
			if strings.HasPrefix(receivedMessageFrame.Text, ":"+aprsHelper.EnsureLength(client.APRSCallAndSSID)+":!") {
				client.SendAck(receivedMessageFrame)
				//strip the prefix
				commandName := strings.ToLower(strings.Split(aprsHelper.ExtractCommand(receivedMessageFrame.Text), " ")[0])
				commandArgs, _ := aprsHelper.ExtractArgs(aprsHelper.ExtractCommand(receivedMessageFrame.Text))
				if commandFunc, exists := commandRegistry[commandName]; exists {
					commandFunc(commandArgs, receivedMessageFrame, client) // Call the corresponding function
				} else {
					fmt.Println("Unknown command:", commandName)
				}
			} else {
				// dont ack acks
				if strings.HasPrefix(receivedMessageFrame.Text, ":"+aprsHelper.EnsureLength(client.CallSign)+":ack") {
					continue
					// dont ack messages not sent to us
				} else if !strings.HasPrefix(receivedMessageFrame.Text, ":"+aprsHelper.EnsureLength(client.CallSign)+":") {
					continue
				} else {
					// if we make it through all that, finally ack the message
					client.SendAck(receivedMessageFrame)
				}
			}
		}
	}
}
