package main

import (
	"io"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/mid"
)

type Recorder struct {
	ports chan<- []string
	c     <-chan Message
}

func NewRecorder() *Recorder {
	var (
		ports = make(chan []string)
		c     = make(chan Message)
	)
	go func() {
		opened := make(map[string]chan struct{})
		for names := range ports {
			required := make(map[string]bool)
			for _, name := range names {
				if _, ok := opened[name]; !ok {
					port, err := mid.OpenIn(midiDriver, -1, name)
					if err != nil {
						continue
					}
					var (
						r, w = io.Pipe()
						msg  = make(chan midi.Message)
						quit = make(chan struct{})
					)
					err = port.SetListener(func(b []byte, deltaMicroseconds int64) {
						w.Write(b)
					})
					if err != nil {
						continue
					}
					go func() {
						mr := midireader.New(r, func(m realtime.Message) {
							msg <- m
						})
						for {
							m, err := mr.Read()
							if err != nil {
								return
							}
							msg <- m
						}
					}()
					go func() {
						for {
							select {
							case m := <-msg:
								c <- Message{port.String(), m}
							case <-quit:
								port.Close()
								return
							}
						}

					}()
					opened[name] = quit
				}
				required[name] = true
			}
			for name := range opened {
				if !required[name] {
					close(opened[name])
					delete(opened, name)
				}
			}
		}
	}()
	return &Recorder{ports, c}
}

func (r *Recorder) Listen(names []string) {
	r.ports <- names
}

func (r *Recorder) C() <-chan Message {
	return r.c
}

type Message struct {
	Port    string
	Message midi.Message
}
