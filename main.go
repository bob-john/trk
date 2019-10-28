package main

import "fmt"

func main() {
	fmt.Println("trk")

	var d midiDriver

	ins, err := d.Ins()
	must(err)

	outs, err := d.Outs()
	must(err)

	for _, port := range ins {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}
	for _, port := range outs {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
