package webaudio

import (
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	c, err := New()
	if c == nil || err != nil {
		t.Fatalf("construct failed: %s %s", c, err)
	}
	c.Close()
}

func TestOscillator(t *testing.T) {
	c, err := New()
	if c == nil || err != nil {
		t.Fatalf("construct failed: %s %s", c, err)
	}
	defer c.Close()
	for i := 1; i <= 4; i++ {
		osc := c.CreateOscillator()
		osc.SetType("sine")
		osc.Frequency().SetValue(float32(440 * i))
		if false {
			osc.Connect(c.Destination())
		} else {
			gain := c.CreateGain()
			gain.Gain().SetValue(0.5 / float32(i))
			osc.Connect(gain)
			gain.Connect(c.Destination())
		}
		osc.Start(0)
		osc.Stop(1.0)
	}
	time.Sleep(time.Second)
}
