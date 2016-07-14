package mock

func Allow(mu Mock) *Allowable {
	return &Allowable{mu: mu}
}

func Â(mu Mock) *Allowable {
	return Allow(mu)
}

type Allowable struct {
	mu   Mock
	last string
}

func (a *Allowable) params(params ...interface{}) []Matcher {
	return nil
}

// On allows the mock to receive a method call with matching parameters and
// return a specific set of values.
func (a *Allowable) On(method string, params ...interface{}) *Allowable {
	calls := a.mu[method]
	calls = append(calls, call{})
	a.mu[method] = calls
	call := &calls[len(calls)-1]
	call.Params = a.params(params...)
	a.last = method
	return a
}

// ToReceive is an alias for Call.
func (a *Allowable) ToReceive(method string, params ...interface{}) *Allowable {
	a.On(method, params...)
	return a
}

// Return specifies what the mock should return when the previously specified
// method is called with matching parameters.
func (a Allowable) Return(results ...interface{}) *Allowable {
	calls := a.mu[a.last]
	if calls == nil || len(calls) < 1 {
		panic("mock: must call On() before calling Return()")
	}
	call := &calls[len(calls)-1]
	if call.Results != nil {
		panic("mock: cannot call Return() twice in a row")
	}
	call.Results = results
	return &a
}
