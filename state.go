package main

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/immesys/bw2bind.v5"
)

// start listening to all of the things and updating the state vector
func (a *Arbiter) subscribe() {
	for _, thing := range a.things {
		thing := thing
		go func(t Thing) {
			c, err := a.client.Subscribe(&bw2bind.SubscribeParams{
				URI: t.ReadURI,
			})
			if err != nil {
				log.Error(errors.Wrapf(err, "Could not subscribe to %s", t.ReadURI))
				return
			}
			for msg := range c {
				a.updateState(t, msg)
			}
		}(thing)
	}

	// every second, save the state and pop onto the stack
	for _ = range time.Tick(1 * time.Second) {
		state := make(map[string]float64)
		a.stateLock.Lock()
		for k, v := range a.state {
			state[k] = v
		}
		a.stateLock.Unlock()
		a.stack <- state
	}
}

func (a *Arbiter) updateState(t Thing, msg *bw2bind.SimpleMessage) {
	var (
		finalValue float64
		value      interface{}
		v          interface{}
	)
	po := msg.POs[0] // get first po
	if obj, ok := po.(bw2bind.MsgPackPayloadObject); !ok {
		// was not msgpack
		log.Debug(po.GetContents())
		bytes := po.GetContents()
		finalValue = float64(bytes[len(bytes)-1]) // grab last value
		goto handleValue
	} else if err := obj.ValueInto(&value); err != nil {
		log.Error(errors.Wrap(err, "Could not unmarshal msgpack"))
	}
	v = t.Getter.Get(value)
	if v_f64, ok := v.(float64); ok {
		finalValue = v_f64
		goto handleValue
	} else if v_i64, ok := v.(int64); ok {
		finalValue = float64(v_i64)
		goto handleValue
	} else if v_u64, ok := v.(uint64); ok {
		finalValue = float64(v_u64)
		goto handleValue
	} else {
		log.Error("don't know value", v)
		return
	}

handleValue:
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	a.state[t.Name] = finalValue
}

func (a *Arbiter) listenActuations() {
	// add things that have arbiter URIs
	for _, thing := range a.things {
		if thing.ArbiterURI == "" {
			continue
		}
		thing := thing
		go func(t Thing) {
			log.Debug("Subscribe to", t.ArbiterURI)
			c, err := a.client.Subscribe(&bw2bind.SubscribeParams{
				URI: t.ArbiterURI,
			})
			if err != nil {
				log.Error(errors.Wrapf(err, "Could not subscribe to %s", t.ReadURI))
				return
			}
			for msg := range c {
				a.generateAction(t, msg)
			}
		}(thing)
	}

}

func (a *Arbiter) generateAction(t Thing, msg *bw2bind.SimpleMessage) {
	log.Debug("got action", t.ArbiterURI)

	po := msg.POs[0]
	var value float64

	if mppo, ok := po.(bw2bind.MsgPackPayloadObject); ok {
		log.Warning("NEED TO IMPLEMENT", mppo)
	} else {
		value = float64(po.GetContents()[0])
	}

	state := map[string]float64{
		t.Name: value,
	}
	action := &publishAction{
		uri:   t.WriteURI,
		state: state,
		ponum: bw2bind.PODFBinaryActuation,
		po:    po,
	}
	log.Debug(action)
	a.actions <- action
}
