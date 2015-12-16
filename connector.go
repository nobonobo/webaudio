// +build !js

package webaudio

import (
	"fmt"
	"sync"
)

type Output interface {
	output(...[]float32)
	connectTo(Input) error
	Disconnect()
	NumberOfOutputs() int
}

type Input interface {
	NumberOfInputs() int
	pull(...[]float32)
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

func (o *OutputImpl) output(buffs ...[]float32) {
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

func (i *InputImpl) pull(buffs ...[]float32) {
	i.RLock()
	defer i.RUnlock()
	for o := range i.sources {
		bs := make([][]float32, len(buffs))
		for ch, buff := range buffs {
			bs[ch] = make([]float32, len(buff))
		}
		o.output(bs...)
		for ch, buff := range buffs {
			for idx, v := range bs[ch] {
				buff[idx] += v
				if buff[idx] > 1.0 {
					buff[idx] = 1.0
				}
				if buff[idx] < -1.0 {
					buff[idx] = -1.0
				}
			}
		}
	}
}
