// +build !js

package webaudio

import (
	"fmt"
	"sync"
)

type Output interface {
	NumberOfOutputs() int
	output() []float32
	connectTo(Input) error
	Disconnect()
}

type Input interface {
	NumberOfInputs() int
	pull([]float32)
	connectFrom(Output)
	disconnectFrom(Output)
}

type OutputImpl struct {
	sync.RWMutex
	numberOfOutputs int
	destinations    map[Input]struct{}
}

type InputImpl struct {
	sync.RWMutex
	numberOfInputs int
	sources        map[Output]struct{}
}

func (o *OutputImpl) output() []float32 {
	panic("must override output method")
}

func (o *OutputImpl) connectTo(i Input) error {
	if i.NumberOfInputs() != o.NumberOfOutputs() {
		return fmt.Errorf("number of channels does not match")
	}
	o.Lock()
	defer o.Unlock()
	o.destinations[i] = struct{}{}
	return nil
}

func (o *OutputImpl) Disconnect() {
	var destinations map[Input]struct{}
	o.Lock()
	o.destinations, destinations = map[Input]struct{}{}, o.destinations
	o.Unlock()
	for i := range destinations {
		i.disconnectFrom(o)
	}
}

func (o *OutputImpl) NumberOfOutputs() int {
	return o.numberOfOutputs
}

func (i *InputImpl) NumberOfInputs() int {
	return i.numberOfInputs
}

func (i *InputImpl) connectFrom(o Output) {
	i.Lock()
	defer i.Unlock()
	i.sources[o] = struct{}{}
}

func (i *InputImpl) disconnectFrom(o Output) {
	i.Lock()
	defer i.Unlock()
	delete(i.sources, o)
}

func (i *InputImpl) pull(values []float32) {
	i.RLock()
	defer i.RUnlock()
	for o := range i.sources {
		outputs := o.output()
		for ch := range outputs {
			values[ch] += outputs[ch%len(outputs)]
		}
	}
}
