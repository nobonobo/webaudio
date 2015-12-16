// +build !js

package webaudio

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"golang.org/x/mobile/exp/audio/al"
)

const (
	CZ  = 2 // bytes/1-sample for al.FormatMono16
	Fmt = al.FormatMono16
)

var SampleRate = float32(44100.0)

type node interface {
	Close()
}

type AudioContext struct {
	sync.RWMutex
	sampleRate    float32
	sampleRateInt int32
	nodes         map[node]struct{}
	source        al.Source
	queue         []al.Buffer
	destination   *AudioDestinationNode
	closed        chan struct{}
	done          chan struct{}
}

func New() (*AudioContext, error) {
	if err := al.OpenDevice(); err != nil {
		return nil, err
	}
	s := al.GenSources(1)
	if code := al.Error(); code != 0 {
		return nil, fmt.Errorf("openal error: %d", code)
	}
	c := new(AudioContext)
	c.sampleRate = SampleRate
	c.sampleRateInt = int32(SampleRate)
	c.nodes = map[node]struct{}{}
	c.source = s[0]
	c.destination = newAudioDestinationNode(c)
	c.closed = make(chan struct{})
	c.done = make(chan struct{})
	qt := 50 * time.Millisecond
	qc := 8
	dt := time.Duration(int(qt) / qc)
	sz := int(c.sampleRate*float32(qt)/float32(time.Second)) / qc
	fmt.Println(sz, qc, dt)
	go func() {
		defer close(c.closed)
		for {
			select {
			case <-c.done:
				return
			default:
				c.proc(sz, qc, dt)
			}
		}
	}()
	return c, nil
}

func (c *AudioContext) proc(sz, qc int, dt time.Duration) {
	n := c.source.BuffersProcessed()
	if n > 0 {
		rm, split := c.queue[:n], c.queue[n:]
		c.queue = split
		c.source.UnqueueBuffers(rm...)
		al.DeleteBuffers(rm...)
	}
	for len(c.queue) < qc {
		n := qc - len(c.queue)
		bs := al.GenBuffers(n)
		for _, b := range bs {
			buf := make([]byte, sz*CZ)
			buff := make([]float32, sz)
			c.destination.pull(buff)
			for idx := 0; idx < sz; idx++ {
				v := int16(float32(32767) * buff[idx])
				binary.LittleEndian.PutUint16(buf[idx*CZ:(idx+1)*CZ], uint16(v))
			}
			b.BufferData(Fmt, buf, c.sampleRateInt)
		}
		c.source.QueueBuffers(bs...)
		c.queue = append(c.queue, bs...)
	}
	if c.source.State() != al.Playing {
		al.PlaySources(c.source)
	}
	time.Sleep(dt)
}

func (c *AudioContext) Close() {
	c.RLock()
	for node := range c.nodes {
		defer node.Close()
	}
	c.RUnlock()
	close(c.done)
	<-c.closed
	//al.CloseDevice()
}

func (c *AudioContext) SampleRate() float32 {
	return c.sampleRate
}

func (c *AudioContext) Destination() *AudioDestinationNode {
	return c.destination
}

func (c *AudioContext) appendNode(n node) {
	c.Lock()
	defer c.Unlock()
	c.nodes[n] = struct{}{}
}

func (c *AudioContext) removeNode(n node) {
	c.Lock()
	defer c.Unlock()
	delete(c.nodes, n)
}

func (c *AudioContext) CreateBuffer(numOfChannels, length, sampleRate int) *AudioBuffer {
	return newAudioBuffer(numOfChannels, length, sampleRate)
}

func (c *AudioContext) CreateOscillator() *OscillatorNode {
	return newOscillatorNode(c)
}

func (c *AudioContext) CreateGain() *GainNode {
	return newGainNode(c)
}
