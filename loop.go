package main

import (
//"time"
)

func (a *Arbiter) Loop() {
	var (
		state     map[string]float64
		action    Action
		gotAction = false
		gotState  = false
	)
	for {
		select {
		case action = <-a.actions:
			gotAction = true
			gotState = false
		case state = <-a.stack:
			gotAction = false
			gotState = true
		}

		if gotAction {
			// check action
			log.Debug("check action", action)
			// merge it
			proposed := action.Proposed()
			for k, v := range proposed {
				state[k] = v
			}
			var ok = REJECTED
			for _, rule := range a.rules {
				if ok = rule.Decide(state); ok == REJECTED {
					log.Warning(rule.Name(), "rejected")
					break
				}
			}
			if ok == APPROVED {
				log.Notice("OK", action)
				action.DoIt(a.client)
			}

		} else if gotState {
			//log.Debug("start rules with state", state)
			for _, rule := range a.rules {
				result := rule.Decide(state)
				if result == REJECTED {
					log.Warning(rule.Name(), "rejected")
				}
				if action := rule.GetAction(); action != nil {
					a.actions <- action
				}
			}
		}
	}
}
