package main

import (
	"github.com/ebarkie/aprs"
	"simpleAPRSbot-go/commands/general"
	"simpleAPRSbot-go/commands/location"
	osuCommands "simpleAPRSbot-go/commands/osu"
	"simpleAPRSbot-go/helpers/aprsHelper"
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
