package main

import (
	"github.com/gomidi/midi"
)

type LaunchpadView interface {
	Update(*Model)
	Handle(*Launchpad, *Model, midi.Message)
	Render(*Launchpad, *Model)
}

type LaunchpadMainView struct {
	active LaunchpadView
}

func NewLaunchpadMainView() *LaunchpadMainView {
	return &LaunchpadMainView{NewLaunchpadSequenceView()}
}

func (v *LaunchpadMainView) Update(model *Model) {
	defer v.ActiveView().Update(model)
}

func (v *LaunchpadMainView) Handle(lp *Launchpad, model *Model, m midi.Message) {
	defer v.ActiveView().Handle(lp, model, m)
	if !lp.IsOn(m) {
		return
	}
	switch lp.Loc(m) {
	case 95:
		v.active = NewLaunchpadSequenceView()
	case 96:
		v.active = NewLaunchpadMuteView()
	case 97:
		v.active = NewLaunchpadPatternView()
	}
}

func (v *LaunchpadMainView) Render(lp *Launchpad, model *Model) {
	v.ActiveView().Render(lp, model)
	lp.ClearModeButtons()
	switch v.ActiveView().(type) {
	case *LaunchpadSequenceView:
		lp.Draw(9, 5, 122)
		lp.Draw(9, 6, 2)
		lp.Draw(9, 7, 2)
	case *LaunchpadMuteView:
		lp.Draw(9, 5, 2)
		lp.Draw(9, 6, 122)
		lp.Draw(9, 7, 2)
	case *LaunchpadPatternView:
		lp.Draw(9, 5, 2)
		lp.Draw(9, 6, 2)
		lp.Draw(9, 7, 122)
	}
}

func (v *LaunchpadMainView) ActiveView() LaunchpadView {
	if v.active == nil {
		v.active = NewLaunchpadSequenceView()
	}
	return v.active
}

type LaunchpadSequenceView struct {
	cursor *Cursor
}

func NewLaunchpadSequenceView() *LaunchpadSequenceView {
	return &LaunchpadSequenceView{NewCursor(3)}
}

func (v *LaunchpadSequenceView) Update(model *Model) {
	v.cursor.Set(8-uint8(model.Cursor()%64)/8, 1+uint8(model.Cursor()%8))
}

func (v *LaunchpadSequenceView) Handle(lp *Launchpad, model *Model, m midi.Message) {
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

func (v *LaunchpadSequenceView) Render(lp *Launchpad, model *Model) {
	lp.ClearNavigationButtons()
	lp.ClearGrid()
	if model.CanDecPage() {
		lp.Draw(9, 1, 2)
	} else {
		lp.Draw(9, 1, 0)
	}
	if model.CanIncPage() {
		lp.Draw(9, 2, 2)
	} else {
		lp.Draw(9, 2, 0)
	}
	v.cursor.Render(lp)
}

type LaunchpadMuteView struct{}

func NewLaunchpadMuteView() *LaunchpadMuteView {
	return &LaunchpadMuteView{}
}

func (v *LaunchpadMuteView) Update(*Model) {

}

func (v *LaunchpadMuteView) Handle(*Launchpad, *Model, midi.Message) {

}

func (v *LaunchpadMuteView) Render(*Launchpad, *Model) {
	lp.ClearNavigationButtons()
	lp.ClearGrid()
}

type LaunchpadPatternView struct{}

func NewLaunchpadPatternView() *LaunchpadPatternView {
	return &LaunchpadPatternView{}
}

func (v *LaunchpadPatternView) Update(*Model) {

}

func (v *LaunchpadPatternView) Handle(*Launchpad, *Model, midi.Message) {

}

func (v *LaunchpadPatternView) Render(lp *Launchpad, model *Model) {
	lp.ClearNavigationButtons()
	lp.ClearGrid()
}
