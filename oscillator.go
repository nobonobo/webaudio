// +build !js

package webaudio

import (
	"fmt"
	"math"
	"sync"

	"golang.org/x/mobile/exp/f32"
)

const (
	Pi = float32(math.Pi)
)

type OscillatorNode struct {
	sync.RWMutex
	waveType    string
	frequency   *AudioParam
	detune      *AudioParam
	phase       float32
	phaseDetune float32
	started     bool
	startDelay  float32
	stopDelay   float32
	generator   func(phase float32) float32
	*AudioNode
	*EventTarget
}

func newOscillatorNode(context *AudioContext) *OscillatorNode {
	o := new(OscillatorNode)
	o.frequency = newAudioParam(440.0)
	o.detune = newAudioParam(0.0)
	o.AudioNode = newAudioNode(context, func() { o.Disconnect() }, 0, 1)
	o.EventTarget = newEventTarget()
	context.appendNode(o.AudioNode)
	o.SetType("sine")
	return o
}

func (o *OscillatorNode) Type() string {
	return o.waveType
}

func (o *OscillatorNode) SetType(t string) error {
	switch t {
	case "sine":
		o.generator = func(phase float32) float32 {
			return f32.Sin(2 * Pi * phase)
		}
	case "square":
		o.generator = func(phase float32) float32 {
			switch {
			case phase < 0.5:
				return 1.0
			default:
				return -1.0
			}
		}
	case "sawtooth":
		o.generator = func(phase float32) float32 {
			switch {
			case phase < 0.5:
				return phase * 2
			default:
				return (phase - 1.0) * 2
			}
		}
	case "triangle":
		o.generator = func(phase float32) float32 {
			switch {
			case phase < 0.25:
				return phase * 4
			case phase < 0.75:
				return 1.0 - (phase-0.25)*4
			default:
				return -1.0 + (phase-0.75)*4
			}
		}
	//case "custom":
	default:
		return fmt.Errorf("unsupported type: %s", t)
	}
	if o.waveType != t {
		o.waveType = t
		o.phase = 0.0
		o.phaseDetune = 0.0
	}
	return nil
}

func (o *OscillatorNode) Frequency() *AudioParam {
	return o.frequency
}

func (o *OscillatorNode) Detune() *AudioParam {
	return o.detune
}

func (o *OscillatorNode) Start(d float32) {
	o.Lock()
	defer o.Unlock()
	if d < 0 || o.started {
		return
	}
	o.startDelay = d
	if d == 0 {
		o.started = true
	}
}

func (o *OscillatorNode) Stop(d float32) {
	o.Lock()
	defer o.Unlock()
	if d < 0 || !o.started {
		return
	}
	o.stopDelay = d
	if d == 0 {
		o.started = false
		o.DispatchEvent(NewEndEvent(o))
	}
}

func (o *OscillatorNode) tick(dt float32) (running, stopped bool) {
	o.Lock()
	defer o.Unlock()
	if o.startDelay > 0 {
		o.startDelay -= dt
		if o.startDelay <= 0 {
			o.startDelay = 0
			o.started = true
		}
	}
	if o.stopDelay > 0 {
		o.stopDelay -= dt
		if o.stopDelay <= 0 {
			o.stopDelay = 0
			o.started = false
			stopped = true
		}
	}
	running = o.started
	return
}

func (o *OscillatorNode) Connect(i Input) error {
	if err := o.connectTo(i); err != nil {
		return err
	}
	i.connectFrom(o)
	return nil
}

func (o *OscillatorNode) output() []float32 {
	dt := 1 / o.context.sampleRate
	freq := o.frequency.output()[0]
	detune := o.detune.output()[0]

	running, stopped := o.tick(dt)
	if stopped {
		o.DispatchEvent(NewEndEvent(o))
	}
	if !running {
		return []float32{0.0}
	}
	value := (o.generator(o.phase) + o.generator(o.phaseDetune)) / 2
	o.phase += freq * dt
	if o.phase > 1.0 {
		o.phase -= 1.0
	}
	o.phaseDetune += detune * dt
	if o.phaseDetune > 1.0 {
		o.phaseDetune -= 1.0
	}
	return []float32{value}
}
