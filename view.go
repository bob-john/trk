package main

// type LaunchpadView interface {
// 	Update(*Model)
// 	Handle(*Launchpad, *Model, midi.Message)
// 	Render(*Launchpad, *Model)
// }

// type LaunchpadMainView struct {
// 	active LaunchpadView
// }

// func NewLaunchpadMainView() *LaunchpadMainView {
// 	return &LaunchpadMainView{NewLaunchpadSequenceView()}
// }

// func (v *LaunchpadMainView) Update(model *Model) {
// 	defer v.ActiveView().Update(model)
// }

// func (v *LaunchpadMainView) Handle(lp *Launchpad, model *Model, m midi.Message) {
// 	defer v.ActiveView().Handle(lp, model, m)
// 	if !lp.IsOn(m) {
// 		return
// 	}
// 	switch lp.Loc(m) {
// 	case 95:
// 		v.active = NewLaunchpadSequenceView()
// 	case 96:
// 		v.active = NewLaunchpadDigitaktView()
// 	case 97:
// 		v.active = NewLaunchpadDigitoneView()
// 	}
// }

// func (v *LaunchpadMainView) Render(lp *Launchpad, model *Model) {
// 	v.ActiveView().Render(lp, model)
// 	lp.ClearModeButtons()
// 	switch v.ActiveView().(type) {
// 	case *LaunchpadSequenceView:
// 		lp.Draw(9, 5, 122)
// 		lp.Draw(9, 6, 2)
// 		lp.Draw(9, 7, 2)
// 	case *LaunchpadDigitaktView:
// 		lp.Draw(9, 5, 2)
// 		lp.Draw(9, 6, 122)
// 		lp.Draw(9, 7, 2)
// 	case *LaunchpadDigitoneView:
// 		lp.Draw(9, 5, 2)
// 		lp.Draw(9, 6, 2)
// 		lp.Draw(9, 7, 122)
// 	}
// }

// func (v *LaunchpadMainView) ActiveView() LaunchpadView {
// 	if v.active == nil {
// 		v.active = NewLaunchpadSequenceView()
// 	}
// 	return v.active
// }

// type LaunchpadSequenceView struct {
// 	cursor *Cursor
// }

// func NewLaunchpadSequenceView() *LaunchpadSequenceView {
// 	return &LaunchpadSequenceView{NewCursor(3, 0)}
// }

// func (v *LaunchpadSequenceView) Update(model *Model) {
// 	v.cursor.Move(8-model.Cursor()%64/8, 1+model.Cursor()%8)
// 	if model.HasMessage(model.Step()) {
// 		v.cursor.SetColor(9)
// 	} else {
// 		v.cursor.SetColor(3)
// 	}
// }

// func (v *LaunchpadSequenceView) Handle(lp *Launchpad, model *Model, m midi.Message) {
// 	if lp.IsOn(m) {
// 		switch lp.Loc(m) {
// 		case 91:
// 			model.DecPage()
// 		case 92:
// 			model.IncPage()
// 		default:
// 			if lp.IsPad(m) {
// 				row, col := lp.Row(m), lp.Col(m)
// 				model.SetCursor(8*(8-row) + col - 1)
// 			}
// 		}
// 	}
// }

// func (v *LaunchpadSequenceView) Render(lp *Launchpad, model *Model) {
// 	lp.ClearNavigationButtons()
// 	lp.ClearGrid()
// 	for cursor := 0; cursor < 64; cursor++ {
// 		step := model.StepForCursor(cursor)
// 		if model.HasMessage(step) {
// 			lp.Draw(8-cursor/8, 1+cursor%8, 5)
// 		} else {
// 			lp.Draw(8-cursor/8, 1+cursor%8, 0)
// 		}
// 	}
// 	if model.CanDecPage() {
// 		lp.Draw(9, 1, 2)
// 	} else {
// 		lp.Draw(9, 1, 0)
// 	}
// 	if model.CanIncPage() {
// 		lp.Draw(9, 2, 2)
// 	} else {
// 		lp.Draw(9, 2, 0)
// 	}
// 	v.cursor.Render(lp)
// }

// type LaunchpadDigitaktView struct{}

// func NewLaunchpadDigitaktView() *LaunchpadDigitaktView {
// 	return &LaunchpadDigitaktView{}
// }

// func (v *LaunchpadDigitaktView) Update(*Model) {}

// func (v *LaunchpadDigitaktView) Handle(*Launchpad, *Model, midi.Message) {}

// func (v *LaunchpadDigitaktView) Render(lp *Launchpad, model *Model) {
// 	lp.ClearNavigationButtons()
// 	lp.ClearGrid()
// }

// type LaunchpadDigitoneView struct{}

// func NewLaunchpadDigitoneView() *LaunchpadDigitoneView {
// 	return &LaunchpadDigitoneView{}
// }

// func (v *LaunchpadDigitoneView) Update(*Model) {}

// func (v *LaunchpadDigitoneView) Handle(lp *Launchpad, model *Model, m midi.Message) {
// 	if !lp.IsOn(m) {
// 		return
// 	}
// 	row, col := lp.Row(m), lp.Col(m)
// 	if col > 8 {
// 		return
// 	}
// 	muted, _ := model.Digitakt.Muted(model.Step(), col-1)
// 	pattern, _ := model.Digitakt.Pattern(model.Step())
// 	switch row {
// 	case 8:
// 		model.Digitakt.SetMuted(model.Step(), col-1, !muted)
// 	case 7:
// 		model.Digitakt.SetPattern(model.Step(), Pattern(pattern).SetTrig(col-1))
// 	case 6:
// 		model.Digitakt.SetPattern(model.Step(), Pattern(pattern).SetTrig(8+col-1))
// 	case 5:
// 		model.Digitakt.SetPattern(model.Step(), Pattern(pattern).SetBank(col-1))
// 	}
// }

// func (v *LaunchpadDigitoneView) Render(lp *Launchpad, model *Model) {
// 	lp.ClearNavigationButtons()
// 	lp.ClearGrid()

// 	// Mute
// 	for i := 0; i < model.Digitakt.VoiceCount(); i++ {
// 		muted, _ := model.Digitakt.Muted(model.Step(), i)
// 		if muted {
// 			lp.Draw(8, 1+i, 0)
// 		} else {
// 			lp.Draw(8, 1+i, 122)
// 		}
// 	}

// 	// Pattern
// 	p, _ := model.Digitakt.Pattern(model.Step())
// 	if p.Trig() < 8 {
// 		lp.Draw(7, 1+p.Trig(), 3)
// 	} else {
// 		lp.Draw(6, 1+p.Trig()%8, 3)
// 	}
// 	lp.Draw(5, 1+p.Bank(), 6)
// }
