package proto // import "docc.io/source/repo/proto"

import (
	"log"
	"os"
)

var NetLogger = log.New(os.Stderr, "net ", log.Ldate|log.Ltime)

type Operations interface {
	// Resolve the latest version.
	Resolve(uri string) (ok bool)

	// Sync to the latest version.
	// Returns whether the version changed.
	Sync(uri string) (ok bool)

	// Archive the state.
	// The directory must be present and writable.
	// Previous entries are discarded.
	Archive(dir string)

	// Extract an archive when present.
	Extract(dir string) (ok bool)

	// Version returns an identifier.
	Version() string
}
