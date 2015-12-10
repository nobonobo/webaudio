// +build !js

package webaudio

import (
	"golang.org/x/mobile/exp/audio/al"
)

type AudioContext struct {
}

func New() (*AudioContext, error) {
	if err := al.OpenDevice(); err != nil {
		return nil, err
	}
	return &AudioContext{}, nil
}

func (c *AudioContext) Close() {
	//al.CloseDevice()
}
