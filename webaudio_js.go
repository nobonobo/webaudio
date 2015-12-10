// +build js

package webaudio

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/util"
)

type AudioContext struct {
	Object *js.Object
	util.EventTarget
}

func New() (*AudioContext, error) {
	var object *js.Object
	if js.Global.Get("AudioContext") != nil {
		object = js.Global.New("AudioContext")
	} else {
		object = js.Global.New("webkitAudioContext")
	}
	if object == nil {
		return nil, fmt.Errorf("not supported")
	}
	return &AudioContext{
		Object:      object,
		EventTarget: util.EventTarget{Object: object},
	}, nil
}

func (c *AudioContext) Close() {
}
