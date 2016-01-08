// +build !js

package webaudio

import "sync"

type GainNode struct {
	sync.RWMutex
	gain   *AudioParam
	_gains []float32
	*AudioNode
	*EventTarget
}

func newGainNode(context *AudioContext) *GainNode {
	g := new(GainNode)
	g.gain = newAudioParam(1.0)
	g.AudioNode = newAudioNode(context, func() { g.Disconnect() }, 1, 1)
	g.EventTarget = newEventTarget()
	context.appendNode(g.AudioNode)
	return g
}

func (g *GainNode) Gain() *AudioParam {
	return g.gain
}

func (g *GainNode) Connect(i Input) error {
	if err := g.connectTo(i); err != nil {
		return err
	}
	i.connectFrom(g)
	return nil
}

func (g *GainNode) output() []float32 {
	gains := g.gain.output()
	buffs := make([]float32, 1)
	g.InputImpl.pull(buffs)
	for ch := range buffs {
		buffs[ch] *= gains[ch%len(gains)]
	}
	return buffs
}
