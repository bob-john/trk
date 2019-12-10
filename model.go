package main

import "trk/track"

type Model struct {
	Track     *track.Track
	Playing   bool
	Recording bool
	Tick      int
	Done      bool
}

func NewModel() *Model {
	return new(Model)
}
