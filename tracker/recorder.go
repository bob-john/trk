package tracker

import (
	"io"
	"log"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midireader"
)

type Recorder struct {
	ports chan<- listenCmd
	c     <-chan Message
}

func NewRecorder(driver mid.Driver) *Recorder {
	var (
		ports = make(chan listenCmd)
		c     = make(chan Message)
	)
	go func() {
		opened := make(map[string]chan struct{})
		for cmd := range ports {
			required := make(map[string]bool)
			for _, name := range cmd.names {
				if _, ok := opened[name]; !ok {
					port, err := mid.OpenIn(driver, -1, name)
					log.Printf("recorder: open %s: %v", port, err)
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
								err := port.Close()
								log.Printf("recorder: close %s: %v", port, err)
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
			if cmd.ack != nil {
				cmd.ack <- struct{}{}
			}
		}
	}()
	return &Recorder{ports, c}
}

func (r *Recorder) Listen(names []string) {
	r.ports <- listenCmd{names, nil}
}

func (r *Recorder) C() <-chan Message {
	return r.c
}

func (r *Recorder) Close() {
	ack := make(chan struct{})
	r.ports <- listenCmd{nil, ack}
	<-ack
}

type Message struct {
	Port    string
	Message midi.Message
}

type listenCmd struct {
	names []string
	ack   chan struct{}
}
