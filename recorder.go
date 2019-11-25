package main

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
					port, err := OpenInput(name)
					if err != nil {
						continue
					}
					quit := make(chan struct{})
					go func() {
						for {
							select {
							case m := <-port.In():
								c <- Message{m, port.Name()}
							case <-quit:
								port.Close()
								return
							}
						}

					}()
					opened[name] = make(chan struct{})
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
