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
	// EachImport enumerates the distinct package names that have been resolved
	// and the import name that was chosen for each. It can be used to generate
	// import statements. The first param of the callback is a package path,
	// the second param is the import name chosen for that path.
	EachImport(func(string, string))
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
	m.pkgNick[pkg] = nick
	m.nickPkg[nick] = pkg
	return fmt.Sprintf("%s.%s", nick, typ)
}

func (m mapResolver) EachImport(cb func(string, string)) {
	for pkg, nick := range m.pkgNick {
		cb(pkg, nick)
	}
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
