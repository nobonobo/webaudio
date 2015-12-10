package webaudio

import "testing"

func TestBasic(t *testing.T) {
	c, err := New()
	if c == nil || err != nil {
		t.Fatalf("construct failed: %s %s", c, err)
	}
	c.Close()
}
