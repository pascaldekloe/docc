package main // import "docc.io/source/cmd/doccd"

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"docc.io/source/proc"
	"docc.io/source/site"
)

var httpAddr = flag.String("http", ":8080", "service address")
var workerCount = flag.Int("workers", 1, "concurreny limit")
var spoolDir = flag.String("spool", proc.WorkRoot, "repository deployment root directory")
var rootDir = flag.String("root", proc.ArchiveRoot, "repository storage root directory")

func main() {
	flag.Parse()

	if err := os.MkdirAll(*spoolDir, 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("%s: unusable spool: %s", os.Args[0], *spoolDir)
	}
	if err := os.MkdirAll(*rootDir, 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("%s: unusable root: %s", os.Args[0], *rootDir)
	}
	proc.WorkRoot = *spoolDir
	proc.ArchiveRoot = *rootDir
	proc.InitWorkers(*workerCount)

	log.Fatal(http.ListenAndServe(*httpAddr, site.Router))
}

func init() {
	signal.Ignore(syscall.SIGHUP)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGKILL)

	go func() {
		sig := <-c
		close(proc.Quit)
		log.Printf("%s: shutdown on signal %q", os.Args[0], sig)

		if sig == syscall.SIGINT {
			<-proc.Done
			os.Exit(0)
		}
	}()
}
