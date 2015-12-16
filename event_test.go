// +build !js

package webaudio

import (
	"fmt"
	"testing"
)

type ev struct {
	typ string
}

func (e *ev) Type() string   { return e.typ }
func (e *ev) String() string { return fmt.Sprintf("Event<%s>", e.typ) }

func TestEvent(t *testing.T) {
	et := newEventTarget()
	onclick := func(e Event) {
		fmt.Println("click!", e.String())
	}
	et.AddEventListener("click", false, onclick)
	et.DispatchEvent(&ev{typ: "click"})
	et.RemoveEventListener("click", false, onclick)
	et.DispatchEvent(&ev{typ: "click"})
}
