package discordgo

import (
	"fmt"
	"testing"
	"time"
)

// TestChannelMessageSend tests the ChannelMessageSend() function. This should not return an error.
func TestSleepCT(t *testing.T) {

	start := time.Now()
	// start the ticker
	s := newSleepCT(20 * time.Millisecond)
	var i int64
	for i = 0; i < 50; i++ {
		s.sleepNext()
	}
	since := time.Since(start)
	fmt.Println("SleepCT after", time.Since(start), "drifts", time.Since(start)-(1*time.Second))
	if since < (980*time.Millisecond) || since > (1020*time.Millisecond) {
		t.Errorf("SleepCT failed timing requirements")
	}
}
