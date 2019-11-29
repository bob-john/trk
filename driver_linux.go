package main

import (
	"gitlab.com/gomidi/midi/mid"
	driver "gitlab.com/gomidi/rtmididrv"
)

func NewDriver() (mid.Driver, error) {
	in, err := rtmidi.NewMIDIInDefault()
	must(err)
	// defer in.Destroy()
	count,err := in.PortCount()
	must(err)
	for i := 0; i < count; i++ {
		name, err := in.PortName(i)
		must(err)
		fmt.Printf("- [%d] %s\n", i,name)
	}
	must(in.OpenPort(0, ""))
	// defer in.Close()
	must(in.IgnoreTypes(false, false, false))
	go func() {
		for  {
			m, t, err := in.Message()
			must(err)
			if len(m) != 0 {
			fmt.Println(m,t)
		}
		}
	}()
	time.Sleep(3*time.Second)
	os.Exit(0)
	return driver.New()
}
