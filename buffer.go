// +build !js

package webaudio

type AudioBuffer struct {
	numOfChannels int
	sampleRate    int
	buffer        [][]float32
}

func newAudioBuffer(numOfChannels, length, sampleRate int) *AudioBuffer {
	b := new(AudioBuffer)
	b.numOfChannels = numOfChannels
	b.sampleRate = sampleRate
	b.buffer = [][]float32{}
	for i := 0; i < numOfChannels; i++ {
		b.buffer = append(b.buffer, make([]float32, length))
	}
	return b
}

func (b *AudioBuffer) ChannelData(n int) []float32 {
	return b.buffer[n]
}
