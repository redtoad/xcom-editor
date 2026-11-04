//go:build !desktop && !dev

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"
)

// OpenURL opens the specified URL in the default browser.
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	var port string
	var host string
	flag.StringVar(&port, "port", "8080", "port to run server on")
	flag.StringVar(&host, "host", "127.0.0.1", "host address to bind to")
	flag.Parse()

	root := flag.Arg(0)
	if root == "" {
		fmt.Fprintln(os.Stderr, "Usage: savegame-editor [-port PORT] [-host HOST] <savegame-root-folder>")
		os.Exit(1)
	}

	log.Printf("Starting savegame editor server...")
	log.Printf("Version %s-%s %s %s", version, commit, date, builtBy)
	log.Printf("Savegame root: %s", root)

	if err := loadRoot(root); err != nil {
		log.Fatal(err)
	}

	r, err := newRouter()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:         host + ":" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Printf("Server running on http://localhost:%s", port)
	log.Println("Opening browser...")
	_ = OpenURL("http://localhost:" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	log.Println("Press Ctrl+C to stop.")
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("Shutting down.")
	os.Exit(0)
}
