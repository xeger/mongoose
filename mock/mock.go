package mock

type allowed struct {
	Params  []Matcher
	Panic   interface{}
	Results []interface{}
}

// State container for mocked behavior.
type Mock map[string][]allowed

// Ø is used to delegate behavior to instances of Mock. Not meant to be called
// directly. If it returns non-nil, then the method call was matched; methods
// that return nothing still return an empty slice.
//
// In contrast, if this method returns nil then the method call was NOT
// matched and the caller should behave accordingly.
//
// If this method receives a nil mock, it panics because Allow() should have
// already initialized the mock or panicked with its own error message.
func Ø(mock Mock, method string, params ...interface{}) []interface{} {
	if mock == nil {
		panic("mock.Ø: cannot delegate behavior to a nil Mock; initialize it before calling")
	}
	calls := mock[method]

	for _, c := range calls {
		if len(c.Params) == len(params) {
			matched := true
			for i, p := range params {
				success, err := c.Params[i].Match(p)
				if err != nil {
					panic(err.Error())
				}
				matched = matched && success
			}
			if matched {
				if c.Panic != nil {
					panic(c.Panic)
				}
				return c.Results
			}
		} else if c.Params == nil {
			return c.Results
		}
	}
	return nil
}
