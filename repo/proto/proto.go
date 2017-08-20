// Package proto provides repository protocol implementations.
package proto // import "docc.io/source/repo/proto"

import (
	"context"
	"log"
	"os"
	"strings"

	"docc.io/source/repo"
)

// NetLogger reports remote connectivity issues.
var NetLogger = log.New(os.Stderr, "net ", log.Ldate|log.Ltime)

// Operations is the work flow abstraction.
type Operations interface {
	// Resolve the latest version.
	Resolve(ctx context.Context) (ok bool)

	// Sync to the latest version.
	// Returns whether the version changed.
	Sync(ctx context.Context) (ok bool)

	// Archive the state.
	// The directory must be present and writable.
	// Previous entries are discarded.
	Archive(dir string)

	// Extract an archive when present.
	Extract(dir string) (ok bool)

	// Version returns an identifier.
	Version() string
}

// MustForName returns a matching implementation for n hosted at path dir.
func MustForName(n repo.Name, dir string) Operations {
	s := string(n)

	switch {
	case strings.HasPrefix(s, "github.com/"):
		return &git{
			URI:  "https://" + s + ".git",
			Root: dir,
		}

	}

	log.Fatalf("%s: unknown host", n)
	panic("unreachable")
}
