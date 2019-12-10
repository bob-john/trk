package main

import "trk/track"

type Model struct {
	Track     *track.Track
	Done      bool
	Playing   bool
	Recording bool
	Tick      int
}

func NewModel() *Model {
	return new(Model)
}
