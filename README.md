# Mongoose

[![Build Status](https://travis-ci.org/xeger/mongoose.png)](https://travis-ci.org/xeger/mongoose)

Mongoose is a tool that parses your Go source code and generates a [mock](https://en.wikipedia.org/wiki/Mock_object) implementation of every [interface](https://gobyexample.com/interfaces) it finds. Mongoose can generate code for a number of mocking packages:

 * [Gomuti](https://github.com/xeger/gomuti)
 * [Testify](https://github.com/stretchr/testify)
 * basic stubs (standalone code; no external package required)

## How to use

Use `go get` to place the `mongoose` binary under `$GOPATH/bin`:

```bash
$ go get github.com/xeger/mongoose
```
### Generate your mocks

Run `mongoose` and pass it the name of one or more directories that contain Go sources. If you are working on package `example/mockitup`:

```bash
$ cd ~/go/src/example/mockitup
$ mongoose .
```

By default, mocks use the [Gomuti](https://github.com/xeger/gomuti) package to record and play back method calls. For information on Testify and other supported mocking toolkits, skip to [Alternative mocking packages](#alternative-mocking-packages) below.

Use the `-r` flag to recurse into all subpackages and generate mocks in parallel.

```bash
$ cd ~/go/src/gigantic/project
$ mongoose -r commerce services util
```

### Write your tests

Consult [the Gomuti documentation](https://github.com/xeger/gomuti) for extensive examples. As a cheat sheet:

```go
import (
  . "github.com/onsi/ginkgo"
  . "github.com/xeger/gomuti"
)

var _ = Describe("stuff", func() {
  It("works", func() {
    // mongoose has generated this type with one method:
    //   Add(l, r int64) int64
    adder := MockAdder{}
    Allow(adder).Call("Add").With(5,5).Return(10)
    Allow(adder).Call("Add").With(10,5).Return(15)
    Allow(adder).Call("Add").With(BeNumerically(">", 2**31-1),Anything()).Panic()

    result := subject.Multiply(3,5))
    Expect(adder).To(HaveCall("Add").Times(2))
    Expect(result).To(Equal(15))    

    Expect(func() {
      subject.Multiply(2**32, 2)
    }).To(Panic())
  })
})
```

## Alternative mocking packages

### Testify

TODO - testify docs

```bash
$ mongoose -mock=testify somepkg
```

Mongoose follows the filesystem conventions of testify's mockery tool; each mock is placed in its own file in the same directory as the mocked interface.

### Plain stubs

WARNING: not implemented yet. Super easy to add but not very useful...

TODO - docs

## How to get help

[Open an issue](https://github.com/xeger/mongoose/issues/new) and explain your problem, steps to reproduce, and your ideal solution (if known).

## How to contribute

Fork the `xeger/mongoose` [repository](https://github.com/xeger/mongoose) on GitHub; make your changes; open a pull request.
