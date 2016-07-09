package parse

import (
	"fmt"
	"regexp"
	"strings"
)

// Namer provides suitable unique names for method parameters or results.
type namer map[string]int

var nonUpper = regexp.MustCompile("[^A-Z]")

func (n namer) Name(pos int, typ Type) string {
	name := strings.ToLower(nonUpper.ReplaceAllString(typ.BareName(), ""))
	_, exists := n[name]
	if exists {
		name = fmt.Sprintf("%s%d", name, pos)
	}
	n[name] = pos
	return name
}
