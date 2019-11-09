package main

import (
	"github.com/gomidi/midi"
)

type LaunchpadView interface {
	Update(*Model)
	Handle(*Launchpad, *Model, midi.Message)
	Render(*Launchpad, *Model)
}

type LaunchpadSessionView struct {
	cursor *Cursor
}

func NewLaunchpadSessionView() *LaunchpadSessionView {
	return &LaunchpadSessionView{NewCursor(3)}
}

func (v *LaunchpadSessionView) Update(model *Model) {
	v.cursor.Set(8-uint8(model.Cursor()%64)/8, 1+uint8(model.Cursor()%8))
}

func (v *LaunchpadSessionView) Handle(lp *Launchpad, model *Model, m midi.Message) {
	if lp.IsOn(m) {
		switch lp.Loc(m) {
		case 91:
			model.DecPage()
		case 92:
			model.IncPage()
		default:
			if lp.IsPad(m) {
				row, col := lp.Row(m), lp.Col(m)
				model.SetCursor(int(8*(8-row) + col - 1))
			}
		}
	}
}

func (v *LaunchpadSessionView) Render(lp *Launchpad, model *Model) {
	if model.CanDecPage() {
		lp.Set(9, 1, 2)
	} else {
		lp.Set(9, 1, 0)
	}
	if model.CanIncPage() {
		lp.Set(9, 2, 2)
	} else {
		lp.Set(9, 2, 0)
	}
	v.cursor.Render(lp)
}
