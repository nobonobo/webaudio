// +build !js

package webaudio

import "sync"

type AudioNode struct {
	sync.RWMutex
	*InputImpl
	*OutputImpl
	once    sync.Once
	context *AudioContext
	release func()
}

func newAudioNode(
	context *AudioContext,
	release func(),
	inputs, outputs int,
) *AudioNode {
	return &AudioNode{
		InputImpl:  &InputImpl{numberOfInputs: inputs, sources: map[Output]struct{}{}},
		OutputImpl: &OutputImpl{numberOfOutputs: outputs, destinations: map[Input]struct{}{}},
		context:    context,
		release:    release,
	}
}

func (n *AudioNode) Close() {
	n.once.Do(func() {
		if n.release != nil {
			n.release()
		}
		n.context.removeNode(n)
	})
}

func (n *AudioNode) Context() *AudioContext {
	return n.context
}
