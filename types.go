// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains custom types

package discordgo

import (
	"time"
)

type sleepCT struct {
	d     time.Duration // desired duration between targets
	t     time.Time     // last time target
	wake  time.Time     // last wake time
	drift int64         // last wake drift microseconds
}

func newSleepCT(d time.Duration) sleepCT {

	s := sleepCT{}

	s.d = d
	s.t = time.Now()

	return s
}

func (s *sleepCT) sleepNext() int64 {

	now := time.Now()

	// if target is zero safety net
	if s.t.IsZero() {
		s.t = now.Add(-s.d)
	}

	// Move forward the sleep target by the duration
	s.t = s.t.Add(s.d)

	// Compute the desired sleep time to reach the target
	d := time.Until(s.t)

	// Sleep
	time.Sleep(d)

	// record the wake time
	s.wake = time.Now()
	s.drift = s.wake.Sub(s.t).Microseconds()

	// return the drift for monitoring purposes
	return s.drift
}
