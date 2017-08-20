package proto // import "docc.io/source/repo/proto"

import (
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

// TestGit is an integration test of the workflow.
func TestGit(t *testing.T) {
	// create repository with one commit
	serveDir, err := ioutil.TempDir("", "git-serve-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(serveDir)

	repoDir := serveDir + "/repo1"
	repoFile := repoDir + "/README"
	{
		if out, err := exec.Command("git", "init", repoDir).CombinedOutput(); err != nil {
			t.Fatalf("Git init: %s - %q", err, out)
		}

		if err := ioutil.WriteFile(repoFile, nil, 0644); err != nil {
			t.Fatal(err)
		}
		add := exec.Command("git", "add", repoFile)
		add.Dir = repoDir
		if out, err := add.CombinedOutput(); err != nil {
			t.Fatalf("Git add: %s - %q", err, out)
		}

		commit := exec.Command("git", "commit", "-m", "first")
		commit.Dir = repoDir
		if out, err := commit.CombinedOutput(); err != nil {
			t.Fatalf("Git first commit: %s - %q", err, out)
		}
	}

	// start server
	rand.Seed(time.Now().UnixNano())
	port := strconv.Itoa(49152 + rand.Intn(65535-49152))

	daemon := exec.Command("git", "daemon", "--port="+port, "--export-all", "--base-path=.", ".")
	daemon.Dir = serveDir
	daemon.Stdout = os.Stdout
	daemon.Stderr = os.Stderr
	if err := daemon.Start(); err != nil {
		t.Fatal("daemon start:", err)
	}
	defer daemon.Process.Signal(os.Interrupt)
	time.Sleep(time.Second)

	// client setup
	var client git
	if client.Root, err = ioutil.TempDir("", "git-work-"); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(client.Root)

	archiveDir, err := ioutil.TempDir("", "git-archive-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(archiveDir)

	// test workflow
	if ok := client.Extract(archiveDir); ok {
		t.Fatal("extracted non-existing archive")
	}

	client.URI = "git://localhost:" + port + "/doesnotexist"
	if ok := client.Resolve(context.Background()); ok {
		t.Fatal("resolved non-existing repository")
	}
	client.URI = "git://localhost:" + port + "/repo1"
	if ok := client.Resolve(context.Background()); !ok {
		t.Fatal("no resolve")
	}
	client.Archive(archiveDir)

	// new client
	client.Root, err = ioutil.TempDir("", "git-work-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(client.Root)

	if ok := client.Extract(archiveDir); !ok {
		t.Fatal("no extraction")
	}
	if ok := client.Sync(context.Background()); ok {
		t.Error("sync without change")
	}

	// update repo
	update := "hello"
	{
		if err := ioutil.WriteFile(repoFile, []byte(update), 0644); err != nil {
			t.Fatal(err)
		}

		commit := exec.Command("git", "commit", "-a", "-m", "second")
		commit.Dir = repoDir
		if out, err := commit.CombinedOutput(); err != nil {
			t.Fatalf("Git second commit: %s - %q", err, out)
		}
	}

	if ok := client.Sync(context.Background()); !ok {
		t.Error("no sync")
	}

	if got, err := ioutil.ReadFile(repoFile); string(got) != update {
		t.Errorf("%s: got %q (%v), want %q", repoFile, got, err, update)
	}
}
