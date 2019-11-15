package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/nsf/termbox-go"
)

const (
	pageSize = 16
)

var (
	doc      *Arrangement
	pen      *Pen
	digitakt *Device
	digitone *Device
	playing  bool
	head     int
)

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Println("usage: trk <path>")
		fmt.Println("trk: invalid command line")
		os.Exit(1)
	}

	doc, err = LoadArrangement(os.Args[1])
	if os.IsNotExist(err) {
		doc = NewArrangement()
	} else if err != nil {
		must(err)
	}

	pen = NewPen(doc)

	err = termbox.Init()
	must(err)
	defer termbox.Close()

	digitakt, _ = OpenDevice("Elektron Digitakt", "Elektron Digitakt")
	digitone, _ = OpenDevice("Elektron Digitone", "Elektron Digitone")

	doc.Row(0).Output(digitakt, digitone)

	var (
		eventC = make(chan termbox.Event)
		midiC  <-chan midi.Message
	)

	if digitakt != nil {
		midiC = digitakt.In()
	} else if digitone != nil {
		midiC = digitone.In()
	}

	go func() {
		for {
			e := termbox.PollEvent()
			if e.Type == termbox.EventInterrupt {
				break
			}
			eventC <- e
		}
	}()

	var (
		done bool
		tick int
	)
	for !done {
		select {
		case e := <-eventC:
			switch e.Type {
			case termbox.EventKey:
				switch e.Key {
				case termbox.KeyCtrlS, termbox.KeyCtrlX:
					doc.WriteFile(os.Args[1])
					done = e.Key == termbox.KeyCtrlX
				}
			}
			if !playing {
				pen.Handle(e, digitakt, digitone)
			}

		case m := <-midiC:
			switch m {
			case realtime.TimingClock:
				if playing {
					tick++
				}
			case realtime.Start:
				playing = true
				tick = 0
				head = 0
			case realtime.Continue:
				playing = true
			case realtime.Stop:
				playing = false
			}
			//TODO Check event on stop and double stop
			// if m != realtime.TimingClock {
			// 	fmt.Println(m)
			// }
		}
		rowLen, _ := strconv.Atoi(doc.Row(head).Len().String())
		if tick == rowLen*6 {
			head++
			tick = 0
		}
		head = clamp(head, 0, doc.RowCount()-1)
		render()
	}
	termbox.Interrupt()
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	SetString(30, 0, fmt.Sprintf("Digitakt: %v", digitakt != nil), termbox.ColorDefault, termbox.ColorDefault)
	SetString(30, 1, fmt.Sprintf("Digitone: %v", digitone != nil), termbox.ColorDefault, termbox.ColorDefault)
	for i := 0; i < pageSize; i++ {
		r := pen.Row() - 8 + i
		if r < 0 || r >= doc.RowCount() {
			continue
		}
		line := doc.Row(r).String()
		fg := termbox.ColorBlue
		if playing && r == head {
			fg = termbox.ColorYellow | termbox.AttrReverse
		}
		SetString(0, i, line, fg, termbox.ColorDefault)
		if !playing && i == 8 {
			SetString(pen.Range().Index, i, pen.Cell().String(), termbox.ColorBlue|termbox.AttrReverse, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func pad(str string, ch rune, total int) string {
	count := total - len(str)
	if count <= 0 {
		return str
	}
	return str + strings.Repeat(string(ch), count)
}

// func colWidth(header string) int {
// 	switch strings.ToUpper(header) {
// 	case "LEN":
// 		return 4
// 	case "DT.P", "DN.P", "P.DT", "P.DN":
// 		return 4 // 3
// 	case "DT.M", "M.DT":
// 		return 8
// 	case "DN.M", "M.DN":
// 		return 4
// 	default:
// 		return len(header)
// 	}
// }

// var (
// 	seq     = []string{}
// 	cur     = Location{}
// 	editing = false

// 	digitakt, digitone *Device
// )

// func main() {
// 	err := termbox.Init()
// 	must(err)
// 	defer termbox.Close()

// 	user, err := user.Current()
// 	must(err)
// 	home := user.HomeDir
// 	appDir := filepath.Join(home, ".trk")
// 	err = os.MkdirAll(appDir, 0700)
// 	must(err)
// 	tmpFilePath := filepath.Join(appDir, "tmp.trk")

// 	seq, err = readFile(tmpFilePath)
// 	if err != nil {
// 		seq = []string{newRow()}
// 	}

// 	digitakt, _ = OpenDevice("Elektron Digitakt", "Elektron Digitakt")
// 	digitone, _ = OpenDevice("Elektron Digitone", "Elektron Digitone")

// 	var (
// 		quit       = make(chan struct{})
// 		transportC = make(chan realtime.Message)
// 	)
// 	defer close(quit)
// 	if digitakt != nil {
// 		defer digitakt.Close()
// 		go listenTransport(digitakt, transportC, quit)
// 	}
// 	if digitone != nil {
// 		defer digitone.Close()
// 	}

// 	go func() {
// 		var (
// 			tick    int
// 			playing bool
// 		)
// 		for {
// 			select {
// 			case m := <-transportC:
// 				switch m {
// 				case realtime.TimingClock:
// 					if playing {
// 						tick++
// 						SetString(30, 0, fmt.Sprintf("%d of %d            ", tick/6, getRowLen()), termbox.ColorDefault, termbox.ColorDefault)
// 						termbox.Flush()
// 						if tick == getRowLen()*6-12 {
// 							if cur.Row+1 < len(seq) {
// 								Row(seq[cur.Row+1]).Play(digitakt, digitone)
// 							}
// 						} else if tick == getRowLen()*6 {
// 							if cur.Row+1 < len(seq) {
// 								cur.Row++
// 								render()
// 							}
// 							tick = 0
// 						}
// 					}
// 				case realtime.Start:
// 					playing = true
// 					tick = 0
// 				case realtime.Continue:
// 					playing = true
// 				case realtime.Stop:
// 					playing = false
// 				}
// 			case <-quit:
// 				return
// 			}
// 		}
// 	}()

// 	var done bool
// 	for !done {
// 		e := termbox.PollEvent()
// 		switch e.Type {
// 		case termbox.EventKey:
// 			switch e.Key {
// 			case termbox.KeyEsc:
// 				done = true

// 			case termbox.KeyArrowUp:
// 				if editing {
// 					getCell(cur).Inc()
// 				} else {
// 					cur.Row--
// 				}
// 			case termbox.KeyPgup:
// 				if editing {
// 					getCell(cur).LargeInc()
// 				} else {
// 					cur.Row -= 16
// 				}
// 			case termbox.KeyHome:
// 				if !editing {
// 					cur.Row = 0
// 				}
// 			case termbox.KeyArrowDown:
// 				if editing {
// 					getCell(cur).Dec()
// 				} else {
// 					cur.Row++
// 				}
// 			case termbox.KeyPgdn:
// 				if editing {
// 					getCell(cur).LargeDec()
// 				} else {
// 					cur.Row += 16
// 				}
// 			case termbox.KeyEnd:
// 				if !editing {
// 					cur.Row = len(seq)
// 				}

// 			case termbox.KeyArrowLeft:
// 				if editing {
// 					cur.Cell--
// 				}
// 			case termbox.KeyArrowRight:
// 				if editing {
// 					cur.Cell++
// 				}

// 			case termbox.KeyEnter:
// 				editing = !editing
// 			case termbox.KeyBackspace, termbox.KeyDelete:
// 				if editing {
// 					getCell(cur).Clear()
// 				} else {
// 					seq = append(seq[:cur.Row], seq[cur.Row+1:]...)
// 					if len(seq) == 0 {
// 						seq = append(seq, newRow())
// 					}
// 				}
// 			case termbox.KeyInsert:
// 				seq = append(seq[:cur.Row+1], append([]string{newRow()}, seq[cur.Row+1:]...)...)
// 				cur.Row++
// 				editing = true
// 			}
// 		}
// 		cur.Row = clamp(cur.Row, 0, len(seq)-1)
// 		cur.Cell = clamp(cur.Cell, 0, 14)
// 		render()
// 		if editing {
// 			writeFile(tmpFilePath)
// 		}
// 		Row(seq[cur.Row]).Play(digitakt, digitone)
// 	}
// }

// func render() {
// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// 	for i := 0; i < 16; i++ {
// 		row := cur.Row - 8 + i
// 		if row < 0 || len(seq) <= row {
// 			continue
// 		}
// 		fg, bg := termbox.ColorBlue, termbox.ColorDefault
// 		if i == 8 && !editing {
// 			fg = fg | termbox.AttrReverse
// 		}
// 		line := seq[row]
// 		SetString(0, i, fmt.Sprintf("%3d %s", 1+row, line), fg, bg)
// 		if i == 8 {
// 			cell := getCell(cur)
// 			SetString(4+cell.Range().Index, i, cell.Range().Substr(line), fg|termbox.AttrReverse, bg)
// 		}
// 	}
// 	termbox.Flush()
// }

// func newRow() string {
// 	return "A01 12345678 A01 1234   16"
// }

// func getCell(c Location) Cell {
// 	line := &seq[c.Row]
// 	switch c.Cell {
// 	case 0:
// 		return NewPatternCell(line, 0, 3)
// 	case 1, 2, 3, 4, 5, 6, 7, 8:
// 		return NewMuteCell(line, 4+c.Cell-1, 1)
// 	case 9:
// 		return NewPatternCell(line, 13, 3)
// 	case 10, 11, 12, 13:
// 		return NewMuteCell(line, 17+c.Cell-10, 1)
// 	case 14:
// 		return NewLenCell(line, 22, 4)
// 	}
// 	return nil
// }

// func getRowLen() int {
// 	return Row(seq[cur.Row]).Len()
// }

// type Row string

// func (r Row) Play(digitakt, digitone *Device) {
// 	if digitakt != nil {
// 		r.Digitakt().Play(digitakt)
// 	}
// 	if digitone != nil {
// 		r.Digitone().Play(digitone)
// 	}
// }

// func (r Row) Len() int {
// 	val, _ := strconv.Atoi(strings.TrimSpace(Range{22, 4}.Substr(string(r))))
// 	return val
// }

// func (r Row) Digitakt() Part {
// 	return Part(string(r)[0:12])
// }

// func (r Row) Digitone() Part {
// 	return Part(string(r)[13:21])
// }

// func (r Row) String() string {
// 	return string(r)
// }

// type Part string

// func (p Part) Play(out *Device) {
// 	if p, ok := p.Pattern(); ok {
// 		out.Write(channel.Channel9.ProgramChange(uint8(p)))
// 	}
// 	p.Mute().Play(out)
// }

// func (p Part) Pattern() (int, bool) {
// 	if strings.ContainsAny(string(p)[0:3], ".") {
// 		return 0, false
// 	}
// 	bank := string(p)[0]
// 	trig, _ := strconv.Atoi(string(p)[1:3])
// 	return int(bank-'A')*16 + trig - 1, true
// }

// func (p Part) Mute() Mute {
// 	return Mute(string(p[4:]))
// }

// type Mute string

// func (m Mute) Play(out *Device) {
// 	c := m.ChannelCount()
// 	for n := 0; n < c; n++ {
// 		var val uint8
// 		switch string(m)[n] {
// 		case '.':
// 			continue
// 		case '+':
// 			val = 0
// 		case '-':
// 			val = 1
// 		}
// 		out.Write(channel.Channel(n).ControlChange(94, val))
// 	}
// }

// func (m Mute) Channel(n int) (bool, bool) {
// 	return string(m)[n] == '-', string(m)[n] != '.'
// }

// func (m Mute) ChannelCount() int {
// 	return len(m)
// }

// func snap(val, grid int) int {
// 	return (val / grid) * grid
// }

// type Location struct {
// 	Row, Cell int
// }

// func writeFile(path string) error {
// 	f, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	for _, line := range seq {
// 		_, err := fmt.Fprintln(f, line)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return f.Close()
// }

// func readFile(path string) ([]string, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	var seq []string
// 	scanner := bufio.NewScanner(f)
// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())
// 		if line != "" {
// 			seq = append(seq, line)
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		return nil, err
// 	}
// 	return seq, f.Close()
// }

// func listenTransport(synth *Device, transportC chan<- realtime.Message, quit <-chan struct{}) {
// 	for {
// 		select {
// 		case m := <-synth.In():
// 			if m, ok := m.(realtime.Message); ok {
// 				transportC <- m
// 			}

// 		case <-quit:
// 			return
// 		}
// 	}
// }

// var (
// 	currentStep        = 0
// 	editing            = false
// 	editor             = &LineEditor{}
// 	seq                = &Seq{}
// 	digitakt, digitone *Device
// )

// func main() {
// 	err := termbox.Init()
// 	must(err)
// 	defer termbox.Close()

// 	user, err := user.Current()
// 	must(err)
// 	home := user.HomeDir

// 	appDir := filepath.Join(home, ".trk")
// 	err = os.MkdirAll(appDir, 0700)
// 	must(err)
// 	tmpFilePath := filepath.Join(appDir, "tmp.trk")

// 	seq.ReadFile(tmpFilePath)

// 	digitakt, _ = OpenDevice("Elektron Digitakt", "Elektron Digitakt")
// 	digitone, _ = OpenDevice("Elektron Digitone", "Elektron Digitone")

// 	var (
// 		quit       = make(chan struct{})
// 		transportC = make(chan realtime.Message)
// 		renderC    = make(chan struct{})
// 		playing    = false
// 		tick       = 0
// 	)
// 	if digitakt != nil {
// 		defer digitakt.Close()
// 		go listenTransport(digitakt, transportC, quit)
// 	}
// 	if digitone != nil {
// 		defer digitone.Close()
// 	}
// 	go func() {
// 		for {
// 			select {
// 			case <-renderC:
// 				render()
// 			case <-quit:
// 				return
// 			}
// 		}
// 	}()
// 	go func() {
// 		for {
// 			select {
// 			case m := <-transportC:
// 				switch m {
// 				case realtime.TimingClock:
// 					if playing && !editing {
// 						tick++
// 						if tick == 12 {
// 							seq.Play(currentStep+1, digitakt, digitone)
// 						} else if tick == 24 {
// 							currentStep++
// 							tick = 0
// 							renderC <- struct{}{}
// 						}
// 					}

// 				case realtime.Start:
// 					playing = true
// 					tick = 0

// 				case realtime.Continue:
// 					playing = true

// 				case realtime.Stop:
// 					playing = false
// 				}

// 			case <-quit:
// 				return
// 			}
// 		}
// 	}()

// 	renderC <- struct{}{}

// 	seq.Play(currentStep, digitakt, digitone)

// 	var done bool
// 	for !done {
// 		e := termbox.PollEvent()
// 		switch e.Type {
// 		case termbox.EventKey:
// 			switch e.Key {
// 			case termbox.KeyEsc:
// 				if editing {
// 					editing = false
// 				} else {
// 					done = true
// 				}

// 			case termbox.KeyArrowUp:
// 				if editing {
// 					editor.ActiveCell().Inc()
// 				} else if currentStep > 0 {
// 					currentStep--
// 					seq.Play(currentStep, digitakt, digitone)
// 				}

// 			case termbox.KeyArrowDown:
// 				if editing {
// 					editor.ActiveCell().Dec()
// 				} else if currentStep < 0xfff {
// 					currentStep++
// 					seq.Play(currentStep, digitakt, digitone)
// 				}

// 			case termbox.KeyDelete, termbox.KeyBackspace:
// 				if editing {
// 					editor.ActiveCell().Clear()
// 				} else {
// 					seq.Insert(seq.emptyLine(currentStep))
// 					seq.WriteFile(tmpFilePath)
// 					seq.Play(currentStep, digitakt, digitone)
// 				}

// 			case termbox.KeyArrowLeft:
// 				if editing {
// 					editor.MoveToPreviousCell()
// 				}

// 			case termbox.KeyArrowRight:
// 				if editing {
// 					editor.MoveToNextCell()
// 				}

// 			case termbox.KeyEnter:
// 				editing = !editing
// 				if editing {
// 					editor.Reset(seq.Line(currentStep), seq.ConsolidatedLine(currentStep))
// 				} else {
// 					seq.Insert(editor.Line())
// 					seq.WriteFile(tmpFilePath)
// 					seq.Play(currentStep, digitakt, digitone)
// 				}

// 			case termbox.KeyPgup:
// 				if editing {
// 					editor.ActiveCell().PageInc()
// 				} else {
// 					currentStep -= 16
// 					if currentStep < 0 {
// 						currentStep = 0
// 					}
// 					seq.Play(currentStep, digitakt, digitone)
// 				}

// 			case termbox.KeyPgdn:
// 				if editing {
// 					editor.ActiveCell().PageDec()
// 				} else {
// 					currentStep += 16
// 					if currentStep > 0xfff {
// 						currentStep = 0xfff
// 					}
// 					seq.Play(currentStep, digitakt, digitone)
// 				}
// 			}
// 		}
// 		renderC <- struct{}{}
// 	}
// 	close(quit)
// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// }

// func listenTransport(synth *Device, transportC chan<- realtime.Message, quit <-chan struct{}) {
// 	for {
// 		select {
// 		case m := <-synth.In():
// 			if m, ok := m.(realtime.Message); ok {
// 				transportC <- m
// 			}

// 		case <-quit:
// 			return
// 		}
// 	}
// }

// func render() {
// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// 	for i := 0; i < 16; i++ {
// 		step := currentStep - 8 + i
// 		if step < 0 || step > 0xfff {
// 			continue
// 		}
// 		fg, bg := termbox.ColorBlue, termbox.ColorDefault
// 		if (step/16)%2 == 1 {
// 			fg = termbox.ColorGreen
// 		}
// 		if i == 8 && !editing {
// 			fg = fg | termbox.AttrReverse
// 		}
// 		line := seq.Line(step)
// 		if i == 8 && editing {
// 			line = editor.Line()
// 		}
// 		SetString(0, i, line, fg, bg)
// 		if i == 8 && editing {
// 			cell := editor.ActiveCell()
// 			SetString(cell.Index(), i, cell.String(), fg|termbox.AttrReverse, bg)
// 		}
// 	}
// 	SetString(30, 0, fmt.Sprintf("DT: %v", digitakt != nil), termbox.ColorDefault, termbox.ColorDefault)
// 	SetString(30, 1, fmt.Sprintf("DN: %v", digitone != nil), termbox.ColorDefault, termbox.ColorDefault)
// 	termbox.Flush()
// }

// // var drv midiDriver
// // var lp *Launchpad
// // var view LaunchpadView
// // var model *Model

// // func main() {
// // 	quit := make(chan struct{})

// // 	err := termbox.Init()
// // 	must(err)
// // 	defer termbox.Close()

// // 	lp, err = ConnectLaunchpad()
// // 	must(err)
// // 	defer lp.Close()
// // 	lp.Reset()

// // 	model = NewModel()
// // 	view = &LaunchpadMainView{}

// // 	render()

// // 	go func() {
// // 		for {
// // 			select {
// // 			case m := <-lp.In():
// // 				view.Handle(lp, model, m)
// // 				render()

// // 			case <-quit:
// // 				return
// // 			}
// // 		}
// // 	}()

// // var done bool
// // for !done {
// // 	e := termbox.PollEvent()
// // 	switch e.Type {
// // 	case termbox.EventKey:
// // 		switch e.Key {
// // 		case termbox.KeyEsc:
// // 			done = true

// // 		case termbox.KeyPgup:
// // 			model.DecPage()
// // 			render()

// // 		case termbox.KeyPgdn:
// // 			model.IncPage()
// // 			render()
// // 		}
// // 	}
// // }
// // close(quit)
// // termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// // }

// // func write(x, y int, s string) {
// // 	for i, c := range s {
// // 		termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
// // 	}
// // }

// // func render() {
// // 	view.Update(model)

// // 	var (
// // 		page = 1 + model.Page()
// // 		bar  = fmt.Sprintf("%d:4", 1+model.Cursor()/16)
// // 		step = 1 + (model.Cursor() % 16)
// // 	)

// // 	write(0, 0, fmt.Sprintf("SEQ %03d PAGE %s TRIG %02d", page, bar, step))
// // 	termbox.Flush()

// // 	view.Render(lp, model)
// // 	lp.Flush()
// // }
