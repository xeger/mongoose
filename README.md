# Mongoose

[![Build Status](https://travis-ci.org/xeger/mongoose.png)](https://travis-ci.org/xeger/mongoose)

Mongoose is a tool that parses your Go source code and generates a [mock](https://en.wikipedia.org/wiki/Mock_object) implementation of every [interface](https://gobyexample.com/interfaces) it finds. Mongoose can generate code for a number of back-end mocking libraries:

 * [Gomuti](https://github.com/xeger/gomuti)
 * [Testify](https://github.com/stretchr/testify)
 * basic stubs (standalone code; no external package required)

## How to use

Use `go get` to place the `mongoose` binary under `#GOPATH/bin`:

```bash
$ go get github.com/xeger/mongoose
```
### Generate your first mocks

Run `mongoose` and pass it the name of a directory that contains Go sources. If you were working on package `example/mockitup`

```
# $ cd ~/go/src/example/mockitup
# $ mongoose
```

By default, Mongoose will generate a file named `mocks.go` in the package directory. Mock types are named after the types they mimic;
if your package had `Widget` and `Shipment` interfaces, then `mocks.go` would contain `MockWidget` and `MockShipment`.

By default, mocks use the [gomuti](https://github.com/xeger/gomuti) package to record and play back method calls. (For information
on Testify and other supported mocking toolkits, skip to "Alternative mocking packages" below.)

In gomuti, each mock is a struct type that exposes only the methods defined on your interface. You can use gomuti's `Allow` method
to program behaviors for any instance of the mock type. To mock the behavior of the `Munge` method, which panics if you pass zero:
```go
import (
  . "github.com/xeger/gomuti"
)

w := &MockWidget{}
Allow(w).ToReceive("Munge").AndPanic("fatal error")
Allow(w).ToReceive("Munge").With(1).AndReturn(true)
```

Every mock can also become  a [stub](http://martinfowler.com/articles/mocksArentStubs.html#TheDifferenceBetweenMocksAndStubs)
and a [spy](https://robots.thoughtbot.com/spy-vs-spy), freeing you from the need to mock "boring" behavior and allowing you
to do [BDD](https://en.wikipedia.org/wiki/Behavior-driven_development) by verifying which methods were actually called, with
which values, and how often.

Rather than expecting specific parameter values, you can use Gomega matchers to define a range or set of permissible values:

```go
Allow(w).ToReceive("Munge").With(BeNumerically(">", 1)).AndReturn(false)
```

For more information about gomuti, its relation to Gomega, and available matchers, consult [the gomuti documentation](https://github.com/xeger/gomuti).

## Alternative mock libraries

### Testify

TODO - docs

### Plain stubs

TODO - docs

## How to get help

[Open an issue](https://github.com/xeger/mongoose/issues/new) and explain your problem, steps to reproduce, and your ideal solution (if known).

## How to contribute

Fork the `xeger/mongoose` [repository](https://github.com/xeger/mongoose) on GitHub; make your changes; open a pull request.
