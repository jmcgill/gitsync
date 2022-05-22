package main

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"io/ioutil"
	"log"
	"path/filepath"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var pid = -1;

func kill() {
	if pid != -1 {
		err := syscall.Kill(pid, syscall.SIGKILL)
		if err != nil {
			panic(err)
		}
	}
}

func pull() {
	binary, lookErr := exec.LookPath("git")
	if lookErr != nil {
		panic(lookErr)
	}
	args := []string{"git", "-C", path, "pull", "origin", "master"}
	pid, err := syscall.ForkExec(binary, args, &syscall.ProcAttr{Files: []uintptr{0, 1, 2}})
	if err != nil {
		panic(err.Error())
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		panic(err.Error())
	}
	_, err = proc.Wait()
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Running done %v", err)
	if err != nil {
		log.Printf("Error in pull")
		panic(err)
	}
}

func run() {
	body, err := ioutil.ReadFile(filepath.Join(path, "Procfile"))
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	args := strings.Split(string(body), " ");
	binary, lookErr := exec.LookPath(args[0])
	if lookErr != nil {
		log.Printf("Error in lookup")
		panic(lookErr)
	}

	env := os.Environ()

	pid, err = syscall.ForkExec(binary, args, &syscall.ProcAttr{Dir: path, Env:env, Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}})
	if err != nil {
		log.Printf("Error in execution")
		panic(err)
	}
}

var path string

func main() {
	path = os.Args[1]

	if _, err := os.Stat(filepath.Join(path, "Procfile")); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("git-sync must be run in a directory containing a Procfile")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// payloadSecret := r.Header.Get("X-Hub-Signature")
		_, err := github.ValidatePayload(r, []byte("A-Cw/2hbF94\\h1()H:Lq22us+(P}"))
		if err != nil {
			log.Printf("Failed to validate webhook: %v", err)
			return
		}
		defer r.Body.Close()

		if github.WebHookType(r) == "push" {
			log.Print("Killing previous process")
			kill()

			// Pull the latest version of the code from the master branch
			log.Print("Pulling latest version of code")
			pull()

			// Run the process
			log.Print("Running...")
			run()
		}

		// event, err := github.ParseWebHook(github.WebHookType(r), payload)
		//if err != nil {
		//	log.Printf("could not parse webhook: err=%s\n", err)
		//	return
		//}
		//
		//switch _ := event.(type) {
		//case *github.PushEvent:
		//	// Kill existing process if running
		//	kill();
		//
		//	// Pull the latest version of the code from the master branch
		//	pull();
		//
		//	// Run the process
		//	run();
		//default:
		//	log.Printf("unknown event type %s\n", github.WebHookType(r))
		//	return
		//}

	})

	log.Printf("Starting sync on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
