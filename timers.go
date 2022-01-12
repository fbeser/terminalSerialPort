package main

import (
	"time"
)

func rawDataTimers() {
	_100miliSec := time.NewTicker(time.Millisecond * 100)
	defer _100miliSec.Stop()
	for {
		select {
		case <-_100miliSec.C:
			rawDataArray()
			if !serialOpen {
				_100miliSec.Stop()
				return
			}
		}
	}
}
