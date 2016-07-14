package mock

type call struct {
	Called  int
	Params  []Matcher
	Results []interface{}
}

// State container for mocked behavior; no public fields or methods. Access
// instances of Mock by calling functions in this package.
type Mock map[string][]call

// Ø is the only method exported by this type; it is used by generated mock
// types to delegate method calls to the underlying Mock object.
//
// Call this method directly,
func (m Mock) Ø(method string, params ...interface{}) []interface{} {
	panic("not implemented")
}
