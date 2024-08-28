package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ebarkie/aprs"
	"log"
	"os"
	"os/signal"
	"simpleAPRSbot-go/helpers/APRS"
	"simpleAPRSbot-go/helpers/api"
	"strings"
	"time"
)

func main() {
	// let's get started! first, build our client
	APRSClient := initAPRSClient()

	// next, crank up our threads!
	// waits for termination so everything shuts down nicely
	// kinda ironic to start up our termination first
	go exitListener()

	// turn on the queue processor thread
	go queueProcessor(APRSClient)

	// and start listening for commands
	go commandListener(APRSClient)
	log.Println("Receiving")

	select {}
}

func exitListener() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down")
	os.Exit(0)
}

func queueProcessor(client *APRS.UserClient) {
	for {
		if len(client.MessageQueue.Queue) <= 0 {
			// do nothing if the queue is empty
			continue
		} else {
			// pull a message out of the queue, and send it
			APRS.SendMessageFrame(client.MessageQueue.Pop())
			// this globally lets us only send a message every x secs. can be turned up or down based on load
			time.Sleep(1 * time.Second)
		}
	}
}

func commandListener(client *APRS.UserClient) {
	for {
		ctx := context.Background()
		fc := aprs.RecvIS(ctx, "rotate.aprs.net:14580", aprs.Addr{Call: client.CallSign, SSID: client.APRSSsid}, client.APRSPassword, "g/"+client.APRSCallAndSSID)
		for receivedMessageFrame := range fc {
			fmt.Println("")
			fmt.Println(receivedMessageFrame)
			fmt.Printf("Received: [%s]\n", receivedMessageFrame.Text)
			fmt.Printf("Expected: [%s]\n", ":"+APRS.EnsureLength(client.APRSCallAndSSID)+":!")
			if strings.HasPrefix(receivedMessageFrame.Text, ":"+APRS.EnsureLength(client.APRSCallAndSSID)+":!") {
				client.SendAck(receivedMessageFrame)
				//strip the prefix
				command, err := APRS.GetCommand(receivedMessageFrame.Text)
				if err != nil {
					log.Println(err)
					continue
				}
				if commandFunc, exists := commandRegistry[command.Name]; exists {
					commandFunc(command.Arguments, receivedMessageFrame, client) // Call the corresponding function
				} else {
					fmt.Println("Unknown command:", command.Name)
				}
			} else {
				// dont ack acks
				if strings.HasPrefix(receivedMessageFrame.Text, ":"+APRS.EnsureLength(client.CallSign)+":ack") {
					continue
					// dont ack messages not sent to us
				} else if !strings.HasPrefix(receivedMessageFrame.Text, ":"+APRS.EnsureLength(client.CallSign)+":") {
					continue
				} else {
					// if we make it through all that, finally ack the message
					client.SendAck(receivedMessageFrame)
				}
			}
		}
	}
}

func initAPRSClient() *APRS.UserClient {
	var aprsCALL = flag.String("APRS_CALL", "N0CALL-0", "N0CALL-0")
	var aprsPass = flag.Int("APRS_PASS", 000000, "00000")
	var APRSFIkey = flag.String("APRS_FI_API_KEY", "", "APRS FI API Key")
	var OpenWeatherMapKey = flag.String("OWM_FI_API_KEY", "", "OpenWeatherMap API Key")
	var osuClientID = flag.Int("OSU_CLIENT_ID", 0, "OSU Client ID")
	var osuClientSecret = flag.String("OSU_CLIENT_SECRET", "", "OSU Client Secret")
	flag.Parse()

	APIClients := api.InitializeAPIClients(api.Keys{
		OsuClientSecret:   *osuClientSecret,
		OsuClientID:       *osuClientID,
		OpenWeatherMapKey: *OpenWeatherMapKey,
		APRSFIkey:         *APRSFIkey,
	})

	return APRS.InitAPRSClient(*aprsCALL, *aprsPass, APIClients)
}
