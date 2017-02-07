package main

import (
	"gopkg.in/immesys/bw2bind.v5"
)

type Action interface {
	Proposed() map[string]float64
	DoIt(*bw2bind.BW2Client)
}

type publishAction struct {
	uri   string
	ponum string
	po    bw2bind.PayloadObject
	state map[string]float64
}

func (a *publishAction) Proposed() map[string]float64 {
	return a.state
}

func (a *publishAction) DoIt(client *bw2bind.BW2Client) {
	client.Publish(&bw2bind.PublishParams{
		URI:            a.uri,
		PayloadObjects: []bw2bind.PayloadObject{a.po},
	})
}
