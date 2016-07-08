package kitchensink

import "net/url"

type FuelCan interface {
	Qty() float32
	Capacity() float32
	Fill(Vehicle) Vehicle
}

// Vehicle can take you places!
type Vehicle interface {
	Range() int
	Drive(dir string, dist float32) url.URL
	Refuel(FuelCan)
}
