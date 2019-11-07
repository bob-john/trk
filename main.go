package main

import (
	"fmt"
	"os"
)

func main() {
	lp, err := ConnectLaunchpad()
	must(err)
	lp.Reset()
	defer func() {
		lp.Clear()
		lp.Update()
		lp.Close()
	}()
	cur := NewCursor(lp, 3)
	cur.MoveTo(1, 1)
	lp.Set(1, 9, 5)
	lp.Update()
	for m := range lp.In() {
		fmt.Println(m)
		if lp.Location(m) == 19 && lp.IsOn(m) {
			return
		}
	}
}

func must(err error) {
	if err != nil {
		fmt.Printf("trk: %v\n", err)
		os.Exit(2)
	}
}
