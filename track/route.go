package track

import "fmt"

type Route struct {
	Input  string
	Output string
	Clock  bool
	ProgCh bool
	Notes  bool
	CC     bool
}

func (r *Route) String() string {
	return fmt.Sprintf("%s -> %s", r.Input, r.Output)
}
