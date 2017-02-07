package main

import (
	"os"
	"sync"
	//"time"

	"github.com/op/go-logging"
	"gopkg.in/immesys/bw2bind.v5"
)

// logger
var log *logging.Logger

func init() {
	log = logging.MustGetLogger("arbiter")
	var format = "%{color}%{level} %{shortfile} %{time:Jan 02 15:04:05} %{color:reset} ▶ %{message}"
	var logBackend = logging.NewLogBackend(os.Stderr, "", 0)
	logBackendLeveled := logging.AddModuleLevel(logBackend)
	logging.SetBackend(logBackendLeveled)
	logging.SetFormatter(logging.MustStringFormatter(format))
}

type Arbiter struct {
	client    *bw2bind.BW2Client
	things    map[string]Thing
	state     map[string]float64
	stateLock sync.Mutex
	rules     []Rule
	stack     chan map[string]float64
	actions   chan Action
}

func main() {
	a := &Arbiter{
		things:  things,
		state:   make(map[string]float64),
		client:  bw2bind.ConnectOrExit(""),
		rules:   generateRules(),
		stack:   make(chan map[string]float64),
		actions: make(chan Action, 1000),
	}

	a.client.SetEntityFromEnvironOrExit()
	a.client.OverrideAutoChainTo(true)

	go a.subscribe()
	go a.listenActuations()

	//go func() {
	//	for _ = range time.Tick(1 * time.Second) {
	//		a.stateLock.Lock()
	//		log.Debugf("%+v", a.state)
	//		a.stateLock.Unlock()
	//	}
	//}()

	a.Loop()
}

// TODO:
// - create archiver object
// - populate it with the things
// - subscribe to all the URIs
// - start generating state vectors
