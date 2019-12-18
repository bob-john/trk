package elektron

import (
	"trk/ui"
)

func Digitakt() *Device {
	return NewDevice(ui.Output("Digitakt"), 8)
}

func Digitone() *Device {
	return NewDevice(ui.Output("Digitone"), 4)
}
