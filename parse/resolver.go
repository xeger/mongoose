package parse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"regexp"
)

// Resolver assigns unique file-local import names ("nicknames") to packages
// that will be imported and used in a Go source file. It works with several
// string representations of packages; to avoid confusion, please keep in mind
// the following terminology:
//    nick, local: one-word, file-local nicknames e.g. "os", "mypkg"
//    typePath: absolute type paths e.g. "/usr/local/src/foo.Widget", "sort.Interface"
//    impPath: relative import paths e.g. "net/url", "github.com/xeger/bar"
type Resolver interface {
	// Resolve transforms an absolute path-and-typename into an imported type with
	// a package-unique nickname. The "local" package name must be given, and if
	// it is the same as the type's package's assigned nickname, the prefix is
	// omitted.
	Resolve(local, typePath string) string

	// Import idempotently adds a package to the resolver's dictionary with a
	// chosen nickname. Returns true if addition was succesful, or if the package
	// is already registered with that name.
	Import(nick, impPath string) bool

	// Imports provides a map of nickname-to-package.
	Imports() map[string]string
}

// NewResolver creates a simple resolver that keeps state with a pair of maps.
func NewResolver() Resolver {
	return &mapResolver{map[string]string{}, map[string]string{}}
}

type mapResolver struct {
	pkgNick map[string]string // path --> nickname
	nickPkg map[string]string // nickname --> path
}

func (m mapResolver) Resolve(local, typePath string) string {
	impPath, typ := m.chew(typePath)

	if impPath == "" {
		return typ // basic type; nothing to do!
	}

	raw_natural := filepath.Base(impPath)
	// Handle import paths that contain disallowed characters such as go-jira, and
	// ldap.v2
	re := regexp.MustCompile(`[-_\.]`)
	natural := re.ReplaceAllString(raw_natural,``)
	if natural == local {
		return typ // type exists locally; no dot prefix
	}

	// First try to use the natural nick for the pkg
	if m.Import(natural, impPath) {
		return fmt.Sprintf("%s.%s", natural, typ)
	}

	// Deal with collisions by appending successively larger integers
	for i := 2; ; i++ {
		nick := fmt.Sprintf("%s%d", natural, i)
		if m.Import(nick, impPath) {
			return fmt.Sprintf("%s.%s", nick, typ)
		}
	}
}

func (m mapResolver) Import(nick, impPath string) bool {
	oldpkg, exists := m.nickPkg[nick]
	if exists {
		return (oldpkg == impPath)
	}

	m.pkgNick[impPath] = nick
	m.nickPkg[nick] = impPath
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

// Extract an impPath and typeName from an absolute typePath. Account for all of
// the following complexities:
//    - basic types (no typePath)
//    - stdlib types (typePath is relative, not absolute)
//    - multiple GOPATH entries
//    - vendored packages
//    - pointer-to and slice-of decorators on typenames
func (m mapResolver) chew(name string) (string, string) {
	gunk := strings.LastIndexAny(name, "[]*")
	//	name = name[gunk+1:]
	if gunk > 0 {
		panic("wtf, gunk in typePath?!?")
	}
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
