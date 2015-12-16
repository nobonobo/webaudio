// +build !js

package webaudio

import (
	"fmt"
	"reflect"
	"sync"
)

// Event ...
type Event interface {
	Type() string
	String() string
}

// EventTarget ...
type EventTarget struct {
	sync.RWMutex
	listeners map[string][]Listener
}

func newEventTarget() *EventTarget {
	return &EventTarget{listeners: map[string][]Listener{}}
}

func (et *EventTarget) remove(typ string, listener Listener) {
	remove := []int{}
	for i, l := range et.listeners[typ] {
		if reflect.DeepEqual(l, listener) {
			remove = append(remove, i)
		}
	}
	for _, i := range remove {
		et.listeners[typ] = append(et.listeners[typ][0:i], et.listeners[typ][i+1:]...)
	}
}

// AddEventListener ...
func (et *EventTarget) AddEventListener(typ string, useCapture bool, cb interface{}) {
	et.Lock()
	defer et.Unlock()
	listener := et.listener(cb)
	if _, ok := et.listeners[typ]; !ok {
		et.listeners[typ] = []Listener{listener}
		return
	}
	et.remove(typ, listener)
	et.listeners[typ] = append(et.listeners[typ], listener)
}

// RemoveEventListener ...
func (et *EventTarget) RemoveEventListener(typ string, useCapture bool, cb interface{}) {
	et.Lock()
	defer et.Unlock()
	listener := et.listener(cb)
	et.remove(typ, listener)
	if len(et.listeners[typ]) == 0 {
		delete(et.listeners, typ)
	}
}

// DispatchEvent ...
func (et *EventTarget) DispatchEvent(e Event) {
	et.RLock()
	for _, l := range et.listeners[e.Type()] {
		defer l.Handle(e)
	}
	et.RUnlock()
}

func (et *EventTarget) listener(cb interface{}) Listener {
	if lstn, ok := cb.(Listener); ok {
		return lstn
	} else {
		return NewListener(cb)
	}
}

// Listener ...
type Listener struct {
	eventType string
	callback  reflect.Value
}

func NewListener(cb interface{}) Listener {
	var e = reflect.TypeOf((*Event)(nil)).Elem()
	var v = reflect.ValueOf(cb)
	var s string
	if v.Kind() != reflect.Func {
		panic("NewListener: must be called with a function as argument")
	}
	if v.IsNil() {
		panic("NewListener: the callback must not be nil")
	}
	if v.Type().NumIn() != 1 {
		panic("NewListener: the callback must accept one argument")
	}
	if t := v.Type().In(0); !t.Implements(e) {
		panic("NewEventListener: the callback must accept a type implementing the Event interface as argument")
	} else if t != e {
		s = reflect.Zero(t).Interface().(Event).Type()
	}

	return Listener{
		eventType: s,
		callback:  v,
	}
}

func (l Listener) Type() string {
	return l.eventType
}

func (l Listener) Handle(e Event) {
	l.callback.Call([]reflect.Value{reflect.ValueOf(e)})
}

type EndEvent struct {
	target interface{}
}

func NewEndEvent(target interface{}) *EndEvent {
	return &EndEvent{target: target}
}

func (e *EndEvent) Type() string   { return "end" }
func (e *EndEvent) String() string { return fmt.Sprintf("Event<end>") }
