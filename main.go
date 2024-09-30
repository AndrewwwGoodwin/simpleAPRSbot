package main

import (
	"context"
	"fmt"
	"github.com/ebarkie/aprs"
	"log"
	"os"
	"os/signal"
	"simpleAPRSbot-go/helpers/APRS"
	"simpleAPRSbot-go/helpers/api"
	"strconv"
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
	var aprsCALL, aprsCallExists = os.LookupEnv("APRS_CALL")
	var aprsPass, aprsPassExists = os.LookupEnv("APRS_PASS")
	var APRSFIkey, aprsFiKeyExists = os.LookupEnv("APRS_FI_API_KEY")
	var OpenWeatherMapKey, OWMKeyExists = os.LookupEnv("OWM_API_KEY")
	var osuClientID, osuClientIdExists = os.LookupEnv("OSU_CLIENT_ID")
	var osuClientSecret, osuClientSecretExists = os.LookupEnv("OSU_CLIENT_SECRET")
	if !aprsCallExists || !aprsPassExists || !aprsFiKeyExists || !OWMKeyExists || !osuClientIdExists || !osuClientSecretExists {
		fmt.Println("APRS_CALL: " + strconv.FormatBool(aprsCallExists))
		fmt.Println("APRS_PASS: " + strconv.FormatBool(aprsPassExists))
		fmt.Println("APRS_FI_API_KEY: " + strconv.FormatBool(aprsFiKeyExists))
		fmt.Println("OWM_API_KEY: " + strconv.FormatBool(OWMKeyExists))
		fmt.Println("OSU_CLIENT_ID: " + strconv.FormatBool(osuClientIdExists))
		fmt.Println("OSU_CLIENT_SECRET: " + strconv.FormatBool(osuClientSecretExists))

		panic("cannot initialize APRS client due to missing environment variables")
	}
	aprsPassConv, err := strconv.Atoi(aprsPass)
	if err != nil {
		log.Println("Error converting APRS pass value to int")
		return nil
	}

	osuClientIDConv, err := strconv.Atoi(osuClientID)
	if err != nil {
		log.Println("Error converting osu client ID to int")
		return nil
	}

	APIClients := api.InitializeAPIClients(api.Keys{
		OsuClientSecret:   osuClientSecret,
		OsuClientID:       osuClientIDConv,
		OpenWeatherMapKey: OpenWeatherMapKey,
		APRSFIkey:         APRSFIkey,
	})

	return APRS.InitAPRSClient(aprsCALL, aprsPassConv, APIClients)
}
