package proc // import "docc.io/source/proc"

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"docc.io/source"
	"docc.io/source/repo"
	"docc.io/source/repo/proto"
)

var (
	// ArchiveRoot is the directory path for repository storage.
	ArchiveRoot = "/srv/docc"

	// WorkRoot is the direcotry path for repository deployment.
	WorkRoot = "/var/spool/docc"

	// DirMode is the creation mask.
	DirMode os.FileMode = 0755
)

type worker int

// pool controls concurrency when not nil.
var pool chan worker

// InitWorkers deploys n processing routines.
func InitWorkers(n int) {
	if pool != nil {
		log.Fatal("worker pool already initiated")
	}
	pool = make(chan worker, n)

	for i := n; i > 0; i-- {
		pool <- worker(i)
	}
	log.Printf("%d workers deployed", n)
}

var (
	// Quit requests the processing to stop.
	Quit = make(chan struct{})

	// Done when the processing stopped.
	Done = make(chan struct{})
)

// ctx is the background context for workers.
var ctx context.Context

func init() {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())

	go func() {
		<-Quit
		cancel()

		n := cap(pool)
		if n > 0 {
			for w := range pool {
				log.Printf("worker%d: drop", w)
				n--
				if n == 0 {
					break
				}
			}
		}

		close(Done)
	}()
}

// Latest resolves the repository interpretation.
func Latest(id repo.Name) []*source.File {
	w := <-pool
	defer func() {
		if w != 0 {
			pool <- w
		}
	}()

	dir, err := ioutil.TempDir(WorkRoot, "worker-")
	if err != nil {
		log.Printf("worker%d abort: %s", w, err)
		return nil
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			log.Printf("worker%d destroy %q: %s", w, id, err)
		}
	}()
	impl := proto.MustForName(id, dir)

	archivePath := ArchiveRoot + "/" + string(id)
	switch {
	case impl.Extract(archivePath):
		impl.Sync(ctx)
	case impl.Resolve(ctx):
		break
	default:
		return nil
	}

	var files []*source.File

	c, err := Repo(dir)
	if err != nil {
		log.Printf("worker%d for %q: %s", w, id, err)
		return nil
	}
	for f := range c {
		files = append(files, f)
	}

	if err := os.MkdirAll(archivePath, DirMode); err != nil && !os.IsExist(err) {
		log.Printf("worker%d archive abort: %s", w, err)
	} else {
		impl.Archive(archivePath)
	}
	return files
}
