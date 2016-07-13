package kitchensink

import "net/url"

type FuelCan interface {
	Qty() float32
	Capacity() float32
	Fill(Vehicle) Vehicle
}

type Wheel interface {
	Diameter() float32
}

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
