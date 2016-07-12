package parse

import (
	"fmt"
	"regexp"
	"strings"
)

// Namer provides suitable unique names for method parameters or results.
type namer map[string]int

var nonUpper = regexp.MustCompile("[^A-Z]")
var nonAlpha = regexp.MustCompile("[^A-Za-z]")

func (n namer) Name(pos int, typ Type) string {
	name := strings.ToLower(nonUpper.ReplaceAllString(typ.BareName(), ""))
	if name == "" {
		prefix := strings.ToLower(nonAlpha.ReplaceAllString(typ.BareName(), ""))[0:1]
		name = fmt.Sprintf("%s%d", prefix, pos)
	}
	_, exists := n[name]
	if exists {
		name = fmt.Sprintf("%s%d", name, pos)
	}
	n[name] = pos
	return name
}
