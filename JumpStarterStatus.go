package main

import ()

type JumpStarterStatus struct {
	RAM     string
	FreeRam string
	//Free    string
	CPU      string
	Dockers  int
	Running  int
	Paused   int
	Stopped  int
	Uptime   string
	Platform string
	Disk     string
	FreeDisk string
}
