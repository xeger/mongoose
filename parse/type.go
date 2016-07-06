package parse

import (
  "go/types"
)

type Type interface {
  Name() string
  // TODO pass lookup table of pkg->importname
  ZeroValue() string
}

type loaderType struct {
  typ types.Type
}

func (lt loaderType) Name() string {
  return lt.typ.String()
}

func (lt loaderType) ZeroValue() string {
  panic("oh noes")
}
