package mock

// Method parameter matcher; convergent with github.com/onsi/gomega.Matcher to
// enable use of Gomega matchers!
type Matcher interface {
	Match(actual interface{}) (success bool, err error)
}
