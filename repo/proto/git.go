package proto // import "docc.io/source/repo/proto"

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	GitResolveTimeout = time.Minute
	GitSyncTimeout    = time.Minute
)

type git struct {
	// Root is the directory path.
	Root string
	// URI is the remote location.
	URI string
}

func (g git) Resolve(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, GitResolveTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "clone", "--quiet", g.URI, ".")
	cmd.Dir = g.Root

	_, err := cmd.Output()
	switch e := err.(type) {
	case nil:
		return true

	case *exec.ExitError:
		NetLogger.Printf("%s: %q: %s - %q", cmd.Dir, strings.Join(cmd.Args, " "), e, e.Stderr)
	default:
		log.Printf("%s: %q: %s", cmd.Dir, strings.Join(cmd.Args, " "), err)
	}
	return false
}

func (g git) Sync(ctx context.Context) bool {
	old := g.Version()
	if old == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, GitResolveTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "pull", "--quiet", g.URI)
	cmd.Dir = g.Root

	_, err := cmd.Output()
	switch e := err.(type) {
	case nil:
		v := g.Version()
		return v != "" && v != old

	case *exec.ExitError:
		NetLogger.Printf("%s: %q: %s - %q", cmd.Dir, strings.Join(cmd.Args, " "), e, e.Stderr)
	default:
		log.Printf("%s: %q: %s", cmd.Dir, strings.Join(cmd.Args, " "), err)
	}
	return false
}

func (g git) Archive(dir string) {
	path := dir + "/git.bundle"

	cmd := exec.Command("git", "bundle", "create", path, "HEAD")
	cmd.Dir = g.Root

	_, err := cmd.Output()
	switch e := err.(type) {
	case nil:
		break
	case *exec.ExitError:
		log.Printf("%s: %q: %s - %q", cmd.Dir, strings.Join(cmd.Args, " "), e, e.Stderr)
	default:
		log.Printf("%s: %q: %s", cmd.Dir, strings.Join(cmd.Args, " "), err)
	}
}

func (g git) Extract(dir string) bool {
	path := dir + "/git.bundle"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	cmd := exec.Command("git", "clone", "--quiet", path, g.Root)
	cmd.Dir = g.Root

	_, err := cmd.Output()
	switch e := err.(type) {
	case nil:
		return true

	case *exec.ExitError:
		log.Printf("%s: %q: %s - %q", cmd.Dir, strings.Join(cmd.Args, " "), e, e.Stderr)
	default:
		log.Printf("%s: %q: %s", cmd.Dir, strings.Join(cmd.Args, " "), err)
	}
	return false
}

func (g git) Version() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = g.Root

	out, err := cmd.Output()
	switch e := err.(type) {
	case nil:
		return strings.TrimSpace(string(out))

	case *exec.ExitError:
		log.Printf("%s: %q: %s - %q", cmd.Dir, strings.Join(cmd.Args, " "), e, e.Stderr)
	default:
		log.Printf("%s: %q: %s", cmd.Dir, strings.Join(cmd.Args, " "), err)
	}
	return ""
}
