package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/redtoad/xcom-editor/savegame"
)

//go:embed templates/*
var Templates embed.FS

type Editor struct {
	config   *Config
	savegame *savegame.Savegame
}

type Config struct {
	Host     string
	Port     string
	GamePath string
}

func NewServer() http.Handler {
	mux := http.NewServeMux()
	var handler http.Handler = mux
	return handler
}

func run(ctx context.Context, w io.Writer, args []string) error {

	// TODO turn args into config values

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.SetOutput(w)

	config := &Config{
		Host:     "127.0.0.1",
		Port:     "8080",
		GamePath: "../../GAME",
	}

	editor := &Editor{
		config: config,
	}
	return editor.Run(ctx)
}

func (e *Editor) Run(ctx context.Context) (err error) {

	handler := http.NewServeMux()
	//handler.Handle("/games/", e.handleGames())
	srv := &http.Server{
		Addr:    net.JoinHostPort(e.config.Host, e.config.Port),
		Handler: handler,
	}

	go func() {
		log.Printf("listening on %s\n", srv.Addr)
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		if serverErr := srv.Shutdown(ctx); serverErr != nil && err == nil {
			err = serverErr
		}
	}()

	wg.Wait()

	return
}

func main2() {
	ctx := context.Background()
	if err := run(ctx, os.Stderr, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
