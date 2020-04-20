package events

import "github.com/g3n/engine/core"

var dispatcher core.IDispatcher

func Init() {
	if dispatcher != nil {
		return
	}

	dispatcher = core.NewDispatcher()
}

func Set(d core.IDispatcher) {
	if dispatcher != nil {
		return
	}

	dispatcher = d
}

func Subscribe(evname string, cb core.Callback) {
	dispatcher.Subscribe(evname, cb)
}

func SubscribeID(evname string, id interface{}, cb core.Callback) {
	dispatcher.SubscribeID(evname, id, cb)
}

func UnsubscribeID(evname string, id interface{}) int {
	return dispatcher.UnsubscribeID(evname, id)
}

func UnsubscribeAllID(id interface{}) int {
	return dispatcher.UnsubscribeAllID(id)
}

func Dispatch(evname string, ev interface{}) int {
	return dispatcher.Dispatch(evname, ev)
}
