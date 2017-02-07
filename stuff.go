package main

import (
	"github.com/gtfierro/ob"
)

type Extractor []ob.Operation

func (e Extractor) Get(v interface{}) interface{} {
	return ob.Eval(e, v)
}

type Thing struct {
	// short memorable name for the thing
	Name string
	// URI where we read the state of the device
	ReadURI string
	// URI where we can write to the device
	WriteURI string
	// Arbiter URI used for intercepting actuations
	// Usually going to be ucberkeley/eop/showroom/arbiter/act/{thing.Name}
	ArbiterURI string
	Getter     Extractor
	// expression to get the state out of an actuation URI. Else assume binary
	ActGetter Extractor
}

// 1 is "on", 0 is "off"
type BinaryThing interface {
	Set(bool)
	TurnOn()
	TurnOff()
	State() bool
}

var _things = []Thing{
	{
		Name:       "building_heat",
		ReadURI:    "ucberkeley/eop/echola/s.powerup.v0/5/i.binact/signal/state",
		WriteURI:   "ucberkeley/eop/echola/s.powerup.v0/5/i.binact/slot/state",
		ArbiterURI: "ucberkeley/eop/showroom/arbiter/act/building_heat",
		Getter:     ob.Parse("Value"),
	},
	{
		Name:       "building_cool",
		ReadURI:    "ucberkeley/eop/echola/s.powerup.v0/6/i.binact/signal/state",
		WriteURI:   "ucberkeley/eop/echola/s.powerup.v0/6/i.binact/slot/state",
		ArbiterURI: "ucberkeley/eop/showroom/arbiter/act/building_cool",
		Getter:     ob.Parse("Value"),
	},
	{
		Name:       "building_fan",
		ReadURI:    "ucberkeley/eop/echola/s.powerup.v0/4/i.binact/signal/state",
		WriteURI:   "ucberkeley/eop/echola/s.powerup.v0/4/i.binact/slot/state",
		ArbiterURI: "ucberkeley/eop/showroom/arbiter/act/building_fan",
		Getter:     ob.Parse("Value"),
	},
	{
		Name:       "building_light",
		ReadURI:    "ucberkeley/eop/echola/s.powerup.v0/8/i.binact/signal/state",
		WriteURI:   "ucberkeley/eop/echola/s.powerup.v0/8/i.binact/slot/state",
		ArbiterURI: "ucberkeley/eop/showroom/arbiter/act/building_light",
		Getter:     ob.Parse("Value"),
	},
	{
		Name:    "hamilton_test_1",
		ReadURI: "amplab/sensors/s.hamilton/00126d070000003c/i.temperature/signal/operative",
		Getter:  ob.Parse("air_temp"),
	},
	{
		Name:    "hamilton_test_2",
		ReadURI: "amplab/sensors/s.hamilton/00126d0700000051/i.temperature/signal/operative",
		Getter:  ob.Parse("air_temp"),
	},
	{
		Name:     "avg_room_temp",
		ReadURI:  "ucberkeley/eop/showroom/arbiter/average_temperature",
		WriteURI: "ucberkeley/eop/showroom/arbiter/average_temperature",
		Getter:   ob.Parse("Value"),
	},
	{
		Name:       "room_setpoint",
		ReadURI:    "ucberkeley/eop/showroom/arbiter/room_setpoint",
		WriteURI:   "ucberkeley/eop/showroom/arbiter/room_setpoint",
		ArbiterURI: "ucberkeley/eop/showroom/arbiter/act/room_setpoint",
		Getter:     ob.Parse("Value"),
	},
}

var things = make(map[string]Thing)

func init() {
	// populate the map
	for _, t := range _things {
		things[t.Name] = t
	}
}
