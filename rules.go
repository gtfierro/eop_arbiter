package main

import (
	"github.com/satori/go.uuid"
	"sync"
	"time"

	"gopkg.in/immesys/bw2bind.v5"
)

var NS = uuid.FromStringOrNil("17cb257a-ecfd-11e6-9163-1002b58053c7")

type Timeseries struct {
	UUID  string
	Time  int64
	Value float64
}

type Decision uint

const (
	APPROVED Decision = iota
	REJECTED
)

type Rule interface {
	// receives a state and returns approve/reject
	Decide(state map[string]float64) Decision
	// either returns nil (in the case of no generated action)
	// or a map of a NEW state vector. The differences between the two
	// state vectors is what gets actuated
	GetAction() Action
	Name() string
}

type averageTemperatureRule struct {
	thing   Thing
	sensors []string
	avg     *float64
	same    bool
	sync.Mutex
}

func (r *averageTemperatureRule) Decide(state map[string]float64) Decision {
	var avg float64
	var num = 0
	for _, sensor := range r.sensors {
		if val, found := state[sensor]; found {
			avg += val
			num += 1
		}
	}
	avg /= float64(num)

	r.Lock()
	defer r.Unlock()
	// don't repeat
	if r.avg != nil && avg == *r.avg {
		r.same = true
		return APPROVED
	}
	// protect against NaN
	if num > 0 {
		r.same = false
		r.avg = &avg
	}

	// always approve
	return APPROVED
}

func (r *averageTemperatureRule) GetAction() Action {
	r.Lock()
	defer r.Unlock()
	if r.same || r.avg == nil {
		return nil
	}
	state := make(map[string]float64)
	state[r.thing.Name] = *r.avg

	action := &publishAction{
		uri:   r.thing.WriteURI,
		state: state,
		ponum: bw2bind.PODFGilesTimeseries,
	}
	ts := Timeseries{
		UUID:  uuid.NewV3(NS, r.thing.Name).String(),
		Time:  time.Now().UnixNano(),
		Value: *r.avg,
	}
	po, err := bw2bind.CreateMsgPackPayloadObject(bw2bind.FromDotForm(action.ponum), ts)
	if err != nil {
		log.Error(err)
		return nil
	}
	action.po = po

	return action
}

func (r *averageTemperatureRule) Name() string {
	return r.thing.Name
}

type roomSetpointRule struct {
	thing          Thing
	avg_temp       string
	tstat_temp     string
	tstat_setpoint string
	setpoint       *float64
	same           bool
	sync.Mutex
}

// TODO: finish this!
func (r *roomSetpointRule) Decide(state map[string]float64) Decision {
	// get avg temperature
	avg_temp := state["avg_room_temp"]

	// compute diff between that and thermostat
	log.Debug(avg_temp)

	r.same = true
	return APPROVED
}

func (r *roomSetpointRule) GetAction() Action {
	r.Lock()
	defer r.Unlock()
	if r.same || r.setpoint == nil {
		return nil
	}
	state := make(map[string]float64)
	state[r.thing.Name] = *r.setpoint

	action := &publishAction{
		uri:   r.thing.WriteURI,
		state: state,
		ponum: bw2bind.PODFGilesTimeseries,
	}
	ts := Timeseries{
		UUID:  uuid.NewV3(NS, r.thing.Name).String(),
		Time:  time.Now().UnixNano(),
		Value: *r.setpoint,
	}
	po, err := bw2bind.CreateMsgPackPayloadObject(bw2bind.FromDotForm(action.ponum), ts)
	if err != nil {
		log.Error(err)
		return nil
	}
	action.po = po

	return action
}

func (r *roomSetpointRule) Name() string {
	return r.thing.Name
}

type simultaneousHeatingCoolingRule struct {
	heater Thing
	cooler Thing
}

func (r *simultaneousHeatingCoolingRule) Decide(state map[string]float64) Decision {
	if state[r.heater.Name] > 0 && state[r.cooler.Name] > 0 {
		return REJECTED
	}
	return APPROVED
}

func (r *simultaneousHeatingCoolingRule) GetAction() Action {
	return nil
}

func (r *simultaneousHeatingCoolingRule) Name() string {
	return "simultaneous heating cooling rule"
}

func generateRules() []Rule {
	var rules []Rule

	rules = append(rules, &averageTemperatureRule{
		thing:   things["avg_room_temp"],
		sensors: []string{"hamilton_test_1", "hamilton_test_2"},
	})
	rules = append(rules, &simultaneousHeatingCoolingRule{
		heater: things["building_heat"],
		cooler: things["building_cool"],
	})
	//rules = append(rules, &roomSetpointRule{
	//	thing:       things["room_setpoint"],
	//	avg_temp:     "avg_room_temp",
	//	tstat_temp:   "venstart uri goes here",
	//	tstat_setpoint: "venstart uri goes here",
	//})

	return rules
}
