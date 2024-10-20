package main

import (
	"github.com/ebarkie/aprs"
	"simpleAPRSbot-go/commands/general"
	"simpleAPRSbot-go/commands/location"
	osuCommands "simpleAPRSbot-go/commands/osu"
	"simpleAPRSbot-go/helpers/APRS"
)

type CommandFunc func(args []string, f aprs.Frame, client *APRS.UserClient)

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
	"loc":      location.Location, // requires aprs.fi
	"location": location.Location, // requires aprs.fi
	"w":        location.Weather,  // requires aprs.fi + weather keys
	"weather":  location.Weather,  // requires aprs.fi + weather keys
	"osu":      osuCommands.Osu,   // requires osu_client_id + osu_client_secret
}
