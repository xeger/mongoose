package parse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Resolver is a stateful component that assigns import names to packages whose
// types appear in type definitions that we process. It remembers any name
// that is assigned to a package and is capable of stripping off path prefixes
// and suffixes.
type Resolver interface {
	// Resolve transforms an absolute path.typename into an imported type with
	// a prefix that is unique to the package. The "local" package name must be
	// given, and if it is the same as the imported package, the prefix is omitted.
	Resolve(local, pathAndType string) string
	// Import adds a package to the resolver's dictionary with a specified name.
	// Returns false if the nickname is already taken.
	Import(nick, pkg string) bool
	// Imports maps nicknames to package paths
	Imports() map[string]string
}

// NewResolver creates a simple resolver that keeps state with a pair of maps.
func NewResolver() Resolver {
	return &mapResolver{map[string]string{}, map[string]string{}}
}

// MapResolver is a simplistic resolver.
type mapResolver struct {
	pkgNick map[string]string // path --> nickname
	nickPkg map[string]string // nickname --> path
}

// Resolve maps "/GOPATH/foo.com/bar.Type" --> ("foo.com/bar", "Type")
func (m mapResolver) Resolve(local, pathAndType string) string {
	absolute, typ := m.chew(pathAndType)
	if absolute == "" {
		return typ // basic type; nothing to do!
	}

	pkg := filepath.Base(absolute)
	if pkg == local {
		return typ // type exists locally; no dot prefix
	}
	nick, ok := m.pkgNick[pkg]
	if ok {
		// nickname already registered for this pkg
		return fmt.Sprintf("%s.%s", nick, typ)
	}

	// allocate a nickname for the new pkg
	nick = filepath.Base(pkg)
	orig := nick
	for i := 2; m.hasNick(nick); i++ {
		// deal with collisions
		nick = fmt.Sprintf("%s%d", orig, i)
	}
	m.Import(nick, pkg)
	return fmt.Sprintf("%s.%s", nick, typ)
}

func (m mapResolver) Import(nick, pkg string) bool {
	if m.hasNick(nick) {
		return false
	}
	m.pkgNick[pkg] = nick
	m.nickPkg[nick] = pkg
	return true
}

func (m mapResolver) Imports() map[string]string {
	return m.nickPkg
}

// Test whether a nickname is taken or reserved.
func (m mapResolver) hasNick(nick string) bool {
	_, ok := m.nickPkg[nick]
	return ok
}

// Turn an absolute path+type into a relative path+type name. Account for
// multiple gopath entries and vendoring.
func (m mapResolver) chew(name string) (string, string) {
	lastDot := strings.LastIndex(name, ".")
	if lastDot < 0 {
		return "", name
	}

	name, typ := name[0:lastDot], name[lastDot+1:]
	vendored := strings.LastIndex(name, "vendor/")
	if vendored > 0 {
		return name[vendored+7:], typ
	}
	for _, el := range strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator)) {
		if strings.Index(name, el) == 0 {
			return name[len(el):], typ
		}
	}
	return name, typ
}
