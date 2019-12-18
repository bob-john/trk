package tracker

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// var (
// 	player *Player
// )

// func init() {
// 	driver, err := rtmididrv.New()
// 	must(err)
// 	player = NewPlayer(driver)
// }

// func Play(e Event) {
// 	must(player.Play(e.Port(), e.Message()))
// }

// type Event interface {
// 	Port() string
// 	Message() midi.Message
// }

// func must(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// type Playable interface {
// 	Play(Device)
// }

// type Device interface {
// 	Play(message midi.Message)
// }
