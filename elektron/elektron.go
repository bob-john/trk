package elektron

import (
	"trk/ui"
)

func Digitakt() *Device {
	return NewDevice(ui.Out("Digitakt"), 8)
}

func Digitone() *Device {
	return NewDevice(ui.Out("Digitone"), 4)
}
