// +build !js

package webaudio

import "sync"

type AudioParam struct {
	sync.RWMutex
	*InputImpl
	defaultValue float32
	value        float32
}

func newAudioParam(defaultValue float32) *AudioParam {
	return &AudioParam{
		InputImpl:    &InputImpl{numberOfInputs: 1, sources: map[Output]struct{}{}},
		defaultValue: defaultValue,
		value:        defaultValue,
	}
}

func (p *AudioParam) DefaultValue() float32 {
	return p.defaultValue
}

func (p *AudioParam) Value() float32 {
	p.RLock()
	defer p.RUnlock()
	return p.value
}

func (p *AudioParam) SetValue(v float32) {
	p.Lock()
	defer p.Unlock()
	p.value = v
}

func (p *AudioParam) output() []float32 {
	p.Lock()
	defer p.Unlock()
	buff := []float32{p.value}
	p.InputImpl.pull(buff)
	p.value = buff[0]
	return buff
}
