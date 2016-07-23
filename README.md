# Mongoose

**WARNING:** Experimental code; see CHANGELOG for planned interface-breaking changes. 

[![Build Status](https://travis-ci.org/xeger/mongoose.png)](https://travis-ci.org/xeger/mongoose)

Mongoose is a tool that parses your Go source code and generates a [mock](https://en.wikipedia.org/wiki/Mock_object) implementation of every [interface](https://gobyexample.com/interfaces) it finds. Mongoose can generate code for a number of back-end mocking packages:

 * [Gomuti](https://github.com/xeger/gomuti)
 * [Testify](https://github.com/stretchr/testify)
 * basic stubs (standalone code; no external package required)

## How to use

Use `go get` to place the `mongoose` binary under `$GOPATH/bin`:

```bash
$ go get github.com/xeger/mongoose
```
### Generate your mocks

Run `mongoose` and pass it the name of a directory that contains Go sources. If you are working on package `example/mockitup`:

```bash
$ cd ~/go/src/example/mockitup
$ mongoose
```

By default, Mongoose will generate a file named `mocks.go` in the package directory. Mock types are named after the types they mimic;
if your package had `Widget` and `Shipment` interfaces, then `mocks.go` would contain `MockWidget` and `MockShipment`.

By default, mocks use the [Gomuti](https://github.com/xeger/gomuti) package to record and play back method calls. (For information on Testify and other supported mocking toolkits, skip to [Alternative mocking packages](#alternative-mocking-packages) below.)

With Gomuti, each mock is a struct type that exposes only the methods defined on your interface. You can use `gomuti.Allow()` to program behaviors for any instance of the mock type. To mock the behavior of the `Munge` method, which panics if you pass zero:
```go
import (
  . "github.com/xeger/gomuti"
)

w := &MockWidget{}
Allow(w).Call("Munge").With(0).Panic("fatal error")
Allow(w).Call("Munge").With(1).Return(true)
```

Rather than expecting specific parameter values, you can use [Gomega](https://github.com/onsi/gomega) matchers to define a range or set of permissible values:

```go
Allow(w).Call("Munge").With(BeNumerically(">", 1)).Return(false)
```

Every mock can also become a [stub](http://martinfowler.com/articles/mocksArentStubs.html#TheDifferenceBetweenMocksAndStubs), freeing you from the need to mock boring behavior:

```go
w := &MockWidget{Stub:true}
w.Munge(1) // returns zero value (false) because no call matches
```

Finally, Gomuti mocks are [spies](https://robots.thoughtbot.com/spy-vs-spy). Custom Gomega matchers allow you to do [BDD](https://en.wikipedia.org/wiki/Behavior-driven_development) with Ginkgo:
```
Expect(w).To( HaveCall("Munge").With( BeNumerically(">", 0) ).Once() )
```

For more information about gomuti, its relation to Gomega, and available matchers, consult [the gomuti documentation](https://github.com/xeger/gomuti).

## Alternative mocking packages

### Testify

TODO - docs

### Plain stubs

TODO - docs

## How to get help

[Open an issue](https://github.com/xeger/mongoose/issues/new) and explain your problem, steps to reproduce, and your ideal solution (if known).

## How to contribute

Fork the `xeger/mongoose` [repository](https://github.com/xeger/mongoose) on GitHub; make your changes; open a pull request.
