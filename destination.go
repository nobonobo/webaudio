// +build !js

package webaudio

type AudioDestinationNode struct {
	*AudioNode
	*EventTarget
}

func newAudioDestinationNode(context *AudioContext) *AudioDestinationNode {
	o := new(AudioDestinationNode)
	o.AudioNode = newAudioNode(context, nil, 1, 0)
	context.appendNode(o.AudioNode)
	return o
}
