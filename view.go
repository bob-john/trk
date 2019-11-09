package main

import (
	"github.com/gomidi/midi"
)

type LaunchpadView interface {
	Update(*Model)
	Handle(*Launchpad, *Model, midi.Message)
	Render(*Launchpad, *Model)
}

type LaunchpadRootView struct {
	current LaunchpadView
}

func NewLaunchpadRootView() *LaunchpadRootView {
	return &LaunchpadRootView{NewLaunchpadSessionView()}
}

func (v *LaunchpadRootView) Handle(lp *Launchpad, model *Model, m midi.Message) {
	defer v.current.Handle(lp, model, m)
	if !lp.IsOn(m) {
		return
	}
	switch lp.Loc(m) {
	case 95:
		v.current = NewLaunchpadSessionView()

	case 96:
		v.current = &LaunchpadDrumsView{}

	case 97:
		v.current = &LaunchpadKeysView{}
	}
}

func (v *LaunchpadRootView) Update(model *Model) {
	v.current.Update(model)
}

func (v *LaunchpadRootView) Render(lp *Launchpad, model *Model) {
	lp.SetHorizontalLine(9, 0, 8, 0)
	v.current.Render(lp, model)
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
	lp.Set(9, 5, 122)
	lp.Set(9, 6, 2)
	lp.Set(9, 7, 2)

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

type LaunchpadDrumsView struct{}

func (v *LaunchpadDrumsView) Update(model *Model) {
}

func (v *LaunchpadDrumsView) Handle(lp *Launchpad, model *Model, m midi.Message) {
}

func (v *LaunchpadDrumsView) Render(lp *Launchpad, model *Model) {
	lp.Set(9, 5, 2)
	lp.Set(9, 6, 122)
	lp.Set(9, 7, 2)
}

type LaunchpadKeysView struct{}

func (v *LaunchpadKeysView) Update(model *Model) {
}

func (v *LaunchpadKeysView) Handle(lp *Launchpad, model *Model, m midi.Message) {
}

func (v *LaunchpadKeysView) Render(lp *Launchpad, model *Model) {
	lp.Set(9, 5, 2)
	lp.Set(9, 6, 2)
	lp.Set(9, 7, 122)
}
