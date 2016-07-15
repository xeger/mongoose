package fixtures

// FuelCan holds fuel.
type FuelCan interface {
	Qty() float32
	Capacity() float32
	Fill(Vehicle) Vehicle
}
