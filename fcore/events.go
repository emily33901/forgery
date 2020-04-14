package fcore

import "github.com/g3n/engine/core"

var e core.IDispatcher

func SetEvents(new core.IDispatcher) {
	e = new
}

func Events() core.IDispatcher {
	if e == nil {
		panic("Someone needs to call SetEvents() first")
	}

	return e
}
