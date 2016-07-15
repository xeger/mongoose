package fixtures

import "net/url"

// Vehicle can take you places!
type Vehicle interface {
	Range() int
	Attach(wheels ...Wheel)
	Wheels() []Wheel
	Drive(dir string, dist float32) url.URL
	Refuel(FuelCan) *FuelCan
	EnterLeave(map[int]string, bool)
	Occupants() map[int]string
}
